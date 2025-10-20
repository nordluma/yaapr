package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
)

type SearchModel struct {
	input textinput.Model
}

func NewSearch() SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Enter anime name"
	ti.Focus()
	ti.CharLimit = 128
	ti.Width = 30

	return SearchModel{input: ti}
}

func (m SearchModel) Init() tea.Cmd     { return textinput.Blink }
func (m SearchModel) IsTransient() bool { return false }

func (m SearchModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			query := m.input.Value()
			if query == "" {
				return m, nil
			}

			loader := NewLoading(fmt.Sprintf("Searching for \"%s\"", query))

			return m, tea.Batch(
				func() tea.Msg { return PushScreenMsg{Screen: loader} },
				// TODO: share client instead of creating it
				searchAnimeCmd(anilist.NewClient(""), query),
			)
		case "esc":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m SearchModel) View() string {
	return fmt.Sprintf(
		"Search Anime:\n\n%s\n\nPress Enter to search",
		m.input.View(),
	)
}

func (m *SearchModel) reset() {
	m.input.SetValue("")
	m.input.Focus()
}

func searchAnimeCmd(client *anilist.AnilistClient, name string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := client.SearchAnime(ctx, name)

		return SearchResponseMsg{Result: res, Err: err}
	}
}
