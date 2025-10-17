package jikan

import (
	"context"
	"testing"
	"time"
)

func TestGetEpisodes(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test with a known anime that has episodes (One Piece - MAL ID: 21)
	showId := 21

	ch, err := GetEpisodes(ctx, showId)
	if err != nil {
		t.Fatalf("GetEpisodes failed: %v", err)
	}

	episodes := make([]Episode, 0)
	timeout := time.After(5 * time.Second)

	for {
		select {
		case ep, ok := <-ch:
			if !ok {
				// Channel closed
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
		t.Logf("Received %d episodes", len(episodes))
		for i, ep := range episodes {
			if i >= 3 { // Just show first 3
				break
			}
			t.Logf("Episode %d: %s", ep.MalID, ep.Title)
		}
	}
}

func TestGetEpisodesInvalidID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test with invalid ID
	showId := -1

	_, err := GetEpisodes(ctx, showId)
	if err == nil {
		t.Error("Expected error for invalid ID")
	}
}
