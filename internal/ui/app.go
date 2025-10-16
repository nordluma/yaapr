package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
}

type AppModel struct {
	stack []Screen
}

func NewApp() AppModel { return AppModel{stack: []Screen{NewStartup()}} }

func (m AppModel) Init() tea.Cmd { return m.current().Init() }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s, cmd := m.current().Update(msg)
	m.stack[len(m.stack)-1] = s

	// TODO: handle esc/enter

	return m, cmd
}

func (m AppModel) View() string { return m.current().View() }

func (m *AppModel) push(s Screen) {
	m.stack = append(m.stack, s)
}

func (m *AppModel) pop() {
	if len(m.stack) > 1 {
		m.stack = m.stack[:len(m.stack)-1]
	}
}

func (m AppModel) current() Screen {
	return m.stack[len(m.stack)-1]
}
