package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
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
