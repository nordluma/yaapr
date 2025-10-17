package ui

import (
	"context"
	"testing"
	"time"

	"github.com/nordluma/yaapr/internal/anilist"
	"github.com/nordluma/yaapr/internal/jikan"
)

func TestAnimeDetailsAndEpisodes(t *testing.T) {
	// Test with a known anime that has episodes (Berserk - AniList ID: 33)
	animeID := 33

	client := anilist.NewClient("")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch anime details
	details, err := client.GetAnimeById(ctx, animeID)
	if err != nil {
		t.Fatalf("Failed to fetch anime details: %v", err)
	}

	t.Logf("Anime: %s (ID: %d, MAL ID: %d)", details.Title.English, details.ID, details.IDMal)

	if details.IDMal == 0 {
		t.Fatal("IDMal is 0")
	}

	// Now try to fetch episodes
	episodeCh, err := jikan.GetEpisodes(ctx, details.IDMal)
	if err != nil {
		t.Fatalf("Failed to get episodes: %v", err)
	}

	episodes := make([]jikan.Episode, 0)
	timeout := time.After(5 * time.Second)

	for {
		select {
		case ep, ok := <-episodeCh:
			if !ok {
				goto done
			}
			episodes = append(episodes, ep)
		case <-timeout:
			t.Fatal("Timeout waiting for episodes")
		}
	}

done:
	if len(episodes) == 0 {
		t.Error("No episodes received")
	} else {
		t.Logf("Received %d episodes for %s", len(episodes), details.Title.English)
	}
}
