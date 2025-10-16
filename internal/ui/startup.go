package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2)
	itemStyle       = lipgloss.NewStyle().PaddingLeft(4)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	helpStyle = list.DefaultStyles().
			HelpStyle.PaddingLeft(4).
			PaddingBottom(1)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(
	w io.Writer,
	m list.Model,
	index int,
	listItem list.Item,
) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type StartupModel struct {
	list list.Model
}

func NewStartup() StartupModel {
	items := []list.Item{
		item("Search Anime"),
		item("Currently Watching"),
		item("Update"),
		item("Settings"),
		item("Quit"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "yaapr"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return StartupModel{list: l}
}

func (m StartupModel) Init() tea.Cmd { return nil }
func (m StartupModel) View() string  { return m.list.View() }

func (m StartupModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.list.SelectedItem().(item)
			switch selected {
			case "Search Anime":
				return m, func() tea.Msg {
					return PushScreenMsg{Screen: NewSearch()}
				}
			case "Currently Watching":
				return m, func() tea.Msg {
					return PushScreenMsg{Screen: NewWatchList()}
				}
			case "Update":
				return m, func() tea.Msg {
					return PushScreenMsg{Screen: NewSync()}
				}
			case "Settings":
				return m, func() tea.Msg {
					return PushScreenMsg{Screen: NewSettings()}
				}
			case "Quit":
				return m, tea.Quit
			}
		case "esc", "q":
			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, cmd
}
