package ui

import tea "github.com/charmbracelet/bubbletea"

type WatchListModel struct{}

func NewWatchList() WatchListModel { return WatchListModel{} }

func (m WatchListModel) Init() tea.Cmd                        { return nil }
func (m WatchListModel) Update(msg tea.Msg) (Screen, tea.Cmd) { return m, nil }

func (m WatchListModel) View() string { return "Currently Watchin [TODO]" }
