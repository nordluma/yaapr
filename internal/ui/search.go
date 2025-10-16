package ui

import tea "github.com/charmbracelet/bubbletea"

type SearchModel struct{}

func NewSearch() SearchModel { return SearchModel{} }

func (m SearchModel) Init() tea.Cmd                        { return nil }
func (m SearchModel) Update(msg tea.Msg) (Screen, tea.Cmd) { return m, nil }

func (m SearchModel) View() string { return "Search Anime [TODO]" }
