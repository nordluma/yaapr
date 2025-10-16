package ui

import tea "github.com/charmbracelet/bubbletea"

type WatchListModel struct{}

func NewWatchList() WatchListModel { return WatchListModel{} }

func (m WatchListModel) Init() tea.Cmd { return nil }
func (m WatchListModel) View() string  { return "Currently Watchin [TODO]" }

func (m WatchListModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}
