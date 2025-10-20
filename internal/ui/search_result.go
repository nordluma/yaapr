package ui

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
)

type SearchResultModel struct {
	results list.Model
}

func NewSearchResults(result []anilist.Anime) SearchResultModel {
	animes := []list.Item{}
	for _, anime := range result {
		animes = append(animes, anime)
	}

	const defaultWidth = 20

	l := list.New(animes, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Search Result"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SearchResultModel{results: l}
}

func (m SearchResultModel) Init() tea.Cmd     { return nil }
func (m SearchResultModel) IsTransient() bool { return false }
func (m SearchResultModel) View() string      { return m.results.View() }

func (m SearchResultModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd
	m.results, cmd = m.results.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected, ok := m.results.SelectedItem().(anilist.Anime)
			if !ok {
				return m, cmd
			}

			loading := NewLoading("Fetching anime details")

			return m, tea.Batch(
				func() tea.Msg { return PushScreenMsg{Screen: loading} },
				// TODO: share client instead of creating it
				fetchAnimeDetailsCmd(anilist.NewClient(""), selected.ID),
			)
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, cmd
}

func fetchAnimeDetailsCmd(client *anilist.AnilistClient, id int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		anime, err := client.GetAnimeById(ctx, id)

		return AnimeDetailsFetchedMsg{Anime: anime, Err: err}
	}
}
