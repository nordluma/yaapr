package ui

import tea "github.com/charmbracelet/bubbletea"

type SyncModel struct{}

func NewSync() SyncModel { return SyncModel{} }

func (m SyncModel) Init() tea.Cmd { return nil }
func (m SyncModel) View() string  { return "Update (Sync with Anilist) [TODO]" }

func (m SyncModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}
