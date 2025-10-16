package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/models"
)

type AnimeDetailsModel struct {
	anime models.Anime
}

func NewAnimeDetails(a models.Anime) AnimeDetailsModel {
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
	return fmt.Sprintf(
		"Anime Details:\n\nTitle: %s\nID: %s\n\n[esc] Back",
		m.anime.Title.English,
		m.anime.ID,
	)
}
