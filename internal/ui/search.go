package ui

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
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
			searchTerm := "berserk"
			loading := NewLoading("Searching for Animes")

			return m, tea.Batch(
				func() tea.Msg { return PushScreenMsg{Screen: loading} },
				searchAnimeCmd(anilist.NewClient(""), searchTerm),
			)
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}

func searchAnimeCmd(client *anilist.Client, name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := client.SearchAnime(ctx, name)

		return SearchResponseMsg{Result: res, Err: err}
	}
}
