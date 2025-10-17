package jikan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const jikanBaseUrl = "https://api.jikan.moe/v4"

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

// Stream all episodes for a show in order, fetching pages concurrently.
//
// The returned channel is closed when all pages have been processed or the
// context is cancelled.
func GetEpisodes(ctx context.Context, showId int) (<-chan Episode, error) {
	outCh := make(chan Episode, 10)
	client := &http.Client{}

	// fetch the first page synchronously to determine how many pages there are
	firstRes, totalPages, err := fetchEpisodePage(ctx, client, showId, 1)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(outCh)

		// send first page immediately
		for _, ep := range firstRes {
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

		// fetch concurrently the remaining pages
		for p := 2; p <= totalPages; p++ {
			wg.Add(1)
			go func(page int) {
				defer wg.Done()

				episodes, _, err := fetchEpisodePage(ctx, client, showId, page)
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
