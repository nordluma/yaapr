package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
	"github.com/nordluma/yaapr/internal/jikan"
)

type AnimeDetailsModel struct {
	anime anilist.AnimeDetails
}

func NewAnimeDetails(a anilist.AnimeDetails) AnimeDetailsModel {
	return AnimeDetailsModel{anime: a}
}

func (m AnimeDetailsModel) Init() tea.Cmd     { return nil }
func (m AnimeDetailsModel) IsTransient() bool { return false }

func (m AnimeDetailsModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			loader := NewLoading(fmt.Sprintf("Fetching episodes for %s", m.anime.Title.Romaji))
			return m, tea.Batch(
				func() tea.Msg { return PushScreenMsg{Screen: loader} },
				fetchEpisodesCmd(m.anime.IDMal),
			)
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}

func (m AnimeDetailsModel) View() string {
	var title string
	if m.anime.Title.English != "" {
		title = m.anime.Title.English
	} else {
		title = m.anime.Title.Romaji
	}

	return fmt.Sprintf(
		"Anime Details:\n\nTitle: %s\nID: %d\nMal ID: %d\nStatus: %s\nGenres: %s\nEpisodes: %d\n\n[esc] / [q] - Back",
		title,
		m.anime.ID,
		m.anime.IDMal,
		m.anime.Status,
		strings.Join(m.anime.Genres, ", "),
		m.anime.Episodes,
	)
}

func fetchEpisodesCmd(showId int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		episodeCh, err := jikan.GetEpisodes(ctx, showId)

		return EpisodesFetchedMsg{EpisodeCh: episodeCh, Err: err}
	}
}
