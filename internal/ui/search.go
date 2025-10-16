package ui

import tea "github.com/charmbracelet/bubbletea"

type SearchModel struct{}

func NewSearch() SearchModel { return SearchModel{} }

func (m SearchModel) Init() tea.Cmd { return nil }
func (m SearchModel) View() string  { return "Search Anime [TODO]" }

func (m SearchModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}
