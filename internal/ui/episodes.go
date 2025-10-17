package ui

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nordluma/yaapr/internal/jikan"
)

type EpisodeModel struct {
	episodes list.Model
	cancel   context.CancelFunc
}

func NewEpisode() EpisodeModel {
	l := list.New([]list.Item{}, itemDelegate{}, 20, listHeight)
	l.Title = "Episodes"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle
	self := EpisodeModel{episodes: l}

	return self
}

func (m EpisodeModel) Init() tea.Cmd     { return nil }
func (m EpisodeModel) IsTransient() bool { return false }
func (m EpisodeModel) View() string      { return m.episodes.View() }

func (m EpisodeModel) Update(msg tea.Msg) (Screen, tea.Cmd) {
	var cmd tea.Cmd
	m.episodes, cmd = m.episodes.Update(msg)

	switch msg := msg.(type) {
	case EpisodeLoadedMsg:
		log.Printf("view received an episode: %s", msg.Episode.Title)
		m.episodes.InsertItem(len(m.episodes.Items()), msg.Episode)

		return m, ListenForEpisodes(msg.EpisodeCh)
	case EpisodeLoadCompleteMsg:
		log.Println("All episodes received")
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			if m.cancel != nil {
				m.cancel()
			}

			return m, func() tea.Msg { return PopScreenMsg{} }
		}
	}

	return m, cmd
}

func (m *EpisodeModel) SetCancel(cancel context.CancelFunc) {
	m.cancel = cancel
}

func ListenForEpisodes(episodeCh <-chan jikan.Episode) tea.Cmd {
	return func() tea.Msg {
		ep, ok := <-episodeCh
		if !ok {
			log.Println("episodes fetched")
			return EpisodeLoadCompleteMsg{}
		}

		log.Printf("got episode: %s\n", ep.Title)
		return EpisodeLoadedMsg{Episode: ep, EpisodeCh: episodeCh}
	}
}
