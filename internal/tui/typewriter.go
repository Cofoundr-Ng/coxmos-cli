package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TypewriterModel struct {
	text     string
	revealed int
	done     bool
	speed    time.Duration
	style    lipgloss.Style
}

func NewTypewriter(text string, speed time.Duration, style lipgloss.Style) TypewriterModel {
	return TypewriterModel{
		text:  text,
		speed: speed,
		style: style,
	}
}

type revealMsg struct{}

func (m TypewriterModel) Init() tea.Cmd {
	return tea.Tick(m.speed, func(t time.Time) tea.Msg {
		return revealMsg{}
	})
}

func (m TypewriterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case revealMsg:
		if m.revealed < len(m.text) {
			m.revealed++
			return m, tea.Tick(m.speed, func(t time.Time) tea.Msg {
				return revealMsg{}
			})
		}
		m.done = true
		return m, tea.Quit
	}
	return m, nil
}

func (m TypewriterModel) View() string {
	return m.style.Render(m.text[:m.revealed])
}

func AnimatedTitle(title string) {
	style := lipgloss.NewStyle().Bold(true).Foreground(Cyan).Padding(0, 1)
	m := NewTypewriter(title, 30*time.Millisecond, style)
	p := tea.NewProgram(m)
	p.Run()
}
