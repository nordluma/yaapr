package ui

import tea "github.com/charmbracelet/bubbletea"

type SyncModel struct{}

func NewSync() SyncModel { return SyncModel{} }

func (m SyncModel) Init() tea.Cmd                        { return nil }
func (m SyncModel) Update(msg tea.Msg) (Screen, tea.Cmd) { return m, nil }

func (m SyncModel) View() string { return "Update (Sync with Anilist) [TODO]" }
