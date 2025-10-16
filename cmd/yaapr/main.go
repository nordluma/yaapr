package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/ui"
)

func main() {
	p := tea.NewProgram(ui.NewApp())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
