package ui

import tea "github.com/charmbracelet/bubbletea"

type SettingsModel struct{}

func NewSettings() SettingsModel { return SettingsModel{} }

func (m SettingsModel) Init() tea.Cmd     { return nil }
func (m SettingsModel) View() string      { return "Settings [TODO]" }
func (m SettingsModel) IsTransient() bool { return false }

func (m SettingsModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, nil
}
