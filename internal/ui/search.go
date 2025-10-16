package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/api"
)

type SearchModel struct{}

func NewSearch() SearchModel { return SearchModel{} }

func (m SearchModel) Init() tea.Cmd     { return nil }
func (m SearchModel) View() string      { return "Search Anime [TODO]" }
func (m SearchModel) IsTransient() bool { return false }

func (m SearchModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := "123"

			loading := NewLoading("Fetching anime details")

			return m, tea.Batch(
				func() tea.Msg { return PushScreenMsg{Screen: loading} },
				fetchAnimeDetailsCmd(selected),
			)
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}

func fetchAnimeDetailsCmd(id string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		anime := api.FetchAnimeDetails(id)

		return AnimeDetailsFetchedMsg{Anime: anime, Err: nil}
	}
}
