package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/models"
)

type SearchModel struct{}

func NewSearch() SearchModel { return SearchModel{} }

func (m SearchModel) Init() tea.Cmd { return nil }
func (m SearchModel) View() string  { return "Search Anime [TODO]" }

func (m SearchModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := models.Anime{
				ID: "123",
				Title: models.AnimeTitle{
					English: "Berserk",
					Romaji:  "Berserk",
				},
			}

			return m, func() tea.Msg {
				return PushScreenMsg{Screen: NewAnimeDetails(selected)}
			}
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}
