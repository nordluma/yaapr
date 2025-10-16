package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type LoadingModel struct {
	spinner spinner.Model
	text    string
}

func NewLoading(text string) LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return LoadingModel{spinner: s, text: text}
}

func (m LoadingModel) Init() tea.Cmd     { return nil }
func (m LoadingModel) IsTransient() bool { return true }

func (m LoadingModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)

	// TODO: remove this after creating an error view
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, cmd
}

func (m LoadingModel) View() string {
	return fmt.Sprintf("%s %s", m.spinner.View(), m.text)
}
