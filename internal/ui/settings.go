package ui

import tea "github.com/charmbracelet/bubbletea"

type SettingsModel struct{}

func NewSettings() SettingsModel { return SettingsModel{} }

func (m SettingsModel) Init() tea.Cmd                        { return nil }
func (m SettingsModel) Update(msg tea.Msg) (Screen, tea.Cmd) { return m, nil }

func (m SettingsModel) View() string { return "Settings [TODO]" }
