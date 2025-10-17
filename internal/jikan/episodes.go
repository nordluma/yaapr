package jikan

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	jikanBaseUrl      = "https://api.jikan.moe/v4"
	maxConcurrency    = 3
	maxRetries        = 3
	rateLimitInterval = time.Second // base rate limit window
	rateLimitBurst    = 3           // up to 3 requests per second
)

var rateLimiter = make(chan struct{}, rateLimitBurst)

func init() {
	ticker := time.NewTicker(rateLimitInterval)
	go func() {
		for range ticker.C {
			for len(rateLimiter) < rateLimitBurst {
				rateLimiter <- struct{}{}
			}
		}
	}()
}

type episodesResponse struct {
	Data       []Episode `json:"data"`
	Pagination struct {
		LastVisiblePage int  `json:"last_visible_page"`
		HasNextPage     bool `json:"has_next_page"`
	} `json:"pagination"`
}

type Episode struct {
	MalID         int    `json:"mal_id"`
	Url           string `json:"url"`
	Title         string `json:"title"`
	TitleJapanese string `json:"title_japanese"`
	TitleRomaji   string `json:"title_romaji"`
	Aired         string `json:"aired"`
	Filler        bool   `json:"filler"`
	Recap         bool   `json:"recap"`
}

func (e Episode) FilterValue() string {
	if e.Title != "" {
		return e.Title
	}

	return e.TitleRomaji
}

// Stream all episodes for a show in order, fetching pages concurrently.
//
// The returned channel is closed when all pages have been processed or the
// context is cancelled.
func GetEpisodes(ctx context.Context, showId int) (<-chan Episode, error) {
	outCh := make(chan Episode, 10)
	client := &http.Client{Timeout: 10 * time.Second}

	// fetch the first page synchronously to determine how many pages there are
	firstPage, totalPages, err := fetchEpisodePageWithRetry(
		ctx, client, showId, 1,
	)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(outCh)

		// send first page immediately
		for _, ep := range firstPage {
			select {
			case <-ctx.Done():
				return
			case outCh <- ep:
			}
		}

		if totalPages <= 1 {
			return
		}

		type pageResult struct {
			page     int
			episodes []Episode
			err      error
		}

		var wg sync.WaitGroup
		results := make(chan pageResult, totalPages-1)
		sem := make(chan struct{}, maxConcurrency)

		// fetch concurrently the remaining pages
		for p := 2; p <= totalPages; p++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			wg.Add(1)
			go func(page int) {
				defer wg.Done()

				// acquire conccurency slot
				select {
				case sem <- struct{}{}:
				case <-ctx.Done():
					return
				}
				defer func() { <-sem }()

				episodes, _, err := fetchEpisodePageWithRetry(
					ctx, client, showId, page,
				)
				select {
				case <-ctx.Done():
					return
				case results <- pageResult{page: page, episodes: episodes, err: err}:
				}
			}(p)
		}

		// close results once all pages have been fetched
		go func() { wg.Wait(); close(results) }()

		next := 2
		buf := make(map[int][]Episode)

		for res := range results {
			if res.err != nil {
				// TODO: should we break or return a partial result
				return
			}

			buf[res.page] = res.episodes

			for {
				if eps, ok := buf[next]; ok {
					for _, ep := range eps {
						select {
						case <-ctx.Done():
							return
						case outCh <- ep:
						}
					}
					delete(buf, next)
					next++
				} else {
					break
				}

				if next > totalPages {
					break
				}
			}
		}
	}()

	return outCh, nil
}

func fetchEpisodePageWithRetry(
	ctx context.Context,
	client *http.Client,
	showID, page int,
) ([]Episode, int, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		case <-rateLimiter: // enforce rate limit
		}

		episodes, totalPages, err := fetchEpisodePage(ctx, client, showID, page)
		if err == nil {
			return episodes, totalPages, nil
		}

		lastErr = err

		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(err, context.Canceled) {
			return nil, 0, err
		}

		var netErr net.Error
		asNetErr := errors.As(err, &netErr)
		retriableErr := netErr.Timeout() || isRateLimitError(err)

		if asNetErr && retriableErr {
			backoff := exponentialBackoff(attempt)
			select {
			case <-ctx.Done():
				return nil, 0, ctx.Err()
			case <-time.After(backoff):
			}

			continue
		}

		break
	}

	return nil, 0, fmt.Errorf(
		"failed after %d retries: %w",
		maxRetries,
		lastErr,
	)
}

// Fetch as single page of episodes and return (episodes, totalPages, error)
func fetchEpisodePage(
	ctx context.Context,
	client *http.Client,
	showId int,
	page int,
) ([]Episode, int, error) {
	url := fmt.Sprintf(
		"%s/anime/%d/episodes?page=%d",
		jikanBaseUrl,
		showId,
		page,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		return nil, 0, fmt.Errorf("rate limited (429)")
	}

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf(
			"non-OK response (%d): %s",
			res.StatusCode,
			body,
		)
	}

	var response episodesResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, response.Pagination.LastVisiblePage, nil
}

func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()

	return strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate")
}

func exponentialBackoff(attempt int) time.Duration {
	// 500ms, 1s, 2s + jitter
	base := 500 * time.Millisecond * time.Duration(
		math.Pow(2, float64(attempt)),
	)
	jitter := time.Duration(rand.Int63n(int64(250 * time.Millisecond)))

	return base + jitter
}
