package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
)

type SearchResultModel struct {
	results list.Model
}

func NewSearchResult(result []anilist.Anime) SearchResultModel {
	animes := []list.Item{}
	for _, anime := range result {
		animes = append(animes, item(anime.Title.English))
	}

	if len(animes) == 0 {
		animes = append(animes, item("No results..."))
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

	return m, cmd
}

func fetchAnimeDetailsCmd(id string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		anime := anilist.FetchAnimeDetails(id)

		return AnimeDetailsFetchedMsg{Anime: anime, Err: nil}
	}
}
