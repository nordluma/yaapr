package ui

import (
	"context"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/anilist"
	"github.com/nordluma/yaapr/internal/jikan"
)

type Screen interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
	IsTransient() bool
}

type AnimeDetailsFetchedMsg struct {
	Anime anilist.AnimeDetails
	Err   error
}

type EpisodesFetchedMsg struct {
	EpisodeCh <-chan jikan.Episode
	Err       error
	Cancel    context.CancelFunc
}

type EpisodeLoadedMsg struct {
	Episode   jikan.Episode
	EpisodeCh <-chan jikan.Episode
}

type EpisodeLoadCompleteMsg struct{}

type SearchResponseMsg struct {
	Result []anilist.Anime
	Err    error
}

type PushScreenMsg struct {
	Screen Screen
}

type PopScreenMsg struct{}

type AppModel struct {
	stack []Screen
}

func NewApp() AppModel { return AppModel{stack: []Screen{NewStartup()}} }

func (m AppModel) Init() tea.Cmd { return m.current().Init() }

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case PushScreenMsg:
		m.push(msg.Screen)
		return m, m.current().Init()
	case PopScreenMsg:
		m.pop()

		return m, nil
	case SearchResponseMsg:
		if msg.Err != nil {
			errorScreen := NewLoading(fmt.Sprintf("Failed to search for anime: %s", msg.Err))
			return m, func() tea.Msg { return PushScreenMsg{Screen: errorScreen} }
		}

		return m, func() tea.Msg {
			return PushScreenMsg{Screen: NewSearchResults(msg.Result)}
		}
	case AnimeDetailsFetchedMsg:
		if msg.Err != nil {
			errorScreen := NewLoading(fmt.Sprintf("Failed to load anime: %s", msg.Err))
			return m, func() tea.Msg { return PushScreenMsg{Screen: errorScreen} }
		}

		return m, func() tea.Msg {
			return PushScreenMsg{Screen: NewAnimeDetails(msg.Anime)}
		}
	case EpisodesFetchedMsg:
		if msg.Err != nil {
			errorScreen := NewLoading(fmt.Sprintf("Failed to load episodes: %s", msg.Err))
			log.Printf("failed to fetch episodes: %s", msg.Err)
			return m, func() tea.Msg { return PushScreenMsg{Screen: errorScreen} }
		}

		episodeScreen := NewEpisode()
		episodeScreen.SetCancel(msg.Cancel)

		return m, tea.Sequence(
			func() tea.Msg { return PushScreenMsg{Screen: episodeScreen} },
			ListenForEpisodes(msg.EpisodeCh),
		)
	case EpisodeLoadedMsg:
		if len(m.stack) > 0 {
			top := m.stack[len(m.stack)-1]
			newScreen, cmd := top.Update(msg)
			m.stack[len(m.stack)-1] = newScreen

			return m, cmd
		}

		return m, nil
	case EpisodeLoadCompleteMsg:
		if len(m.stack) > 0 {
			top := m.stack[len(m.stack)-1]
			newScreen, cmd := top.Update(msg)
			m.stack[len(m.stack)-1] = newScreen

			return m, cmd
		}

		return m, nil
	}

	// pass all other messages to the current screen
	if len(m.stack) > 0 {
		top := m.stack[len(m.stack)-1]
		newScreen, cmd := top.Update(msg)
		m.stack[len(m.stack)-1] = newScreen

		return m, cmd
	}

	// TODO: handle esc

	return m, nil
}

func (m AppModel) View() string { return m.current().View() }

func (m *AppModel) push(s Screen) {
	m.stack = append(m.stack, s)
}

func (m *AppModel) pop() {
	if len(m.stack) > 1 {
		// pop current screen
		m.stack = m.stack[:len(m.stack)-1]

		// pop transient screens
		for len(m.stack) > 0 && m.stack[len(m.stack)-1].IsTransient() {
			m.stack = m.stack[:len(m.stack)-1]
		}
	}
}

func (m AppModel) current() Screen {
	return m.stack[len(m.stack)-1]
}
