package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SpinnerModel struct {
	title    string
	done     bool
	err      error
	elapsed  time.Duration
	start    time.Time
	quitting bool
}

func NewSpinner(title string) SpinnerModel {
	return SpinnerModel{
		title: title,
		start: time.Now(),
	}
}

type tickMsg time.Time

func (m SpinnerModel) Init() tea.Cmd {
	return tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.elapsed = time.Since(m.start)
		if m.quitting {
			return m, tea.Quit
		}
		return m, tea.Tick(60*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case error:
		m.err = msg
		m.done = true
		m.quitting = true
		return m, tea.Quit
	case string:
		if msg == "done" {
			m.done = true
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m SpinnerModel) View() string {
	if m.quitting {
		if m.err != nil {
			return lipgloss.NewStyle().Foreground(Red).Render("✗") + " " + m.title + "\n" +
				lipgloss.NewStyle().Foreground(Red).Italic(true).Render("  "+m.err.Error())
		}
		return lipgloss.NewStyle().Foreground(Green).Render("✓") + " " + m.title +
			lipgloss.NewStyle().Foreground(Gray).Italic(true).Render(fmt.Sprintf(" (%.1fs)", m.elapsed.Seconds()))
	}
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	f := frames[int(m.elapsed.Milliseconds()/60)%len(frames)]
	return lipgloss.NewStyle().Foreground(Cyan).Render(f) + " " + m.title
}

type ProgressModel struct {
	title    string
	progress float64
	done     bool
	start    time.Time
	elapsed  time.Duration
	quitting bool
}

func NewProgress(title string) ProgressModel {
	return ProgressModel{title: title, start: time.Now()}
}

func (m ProgressModel) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.elapsed = time.Since(m.start)
		if m.quitting {
			return m, tea.Quit
		}
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case float64:
		m.progress = msg
		if msg >= 1.0 {
			m.done = true
			m.quitting = true
			return m, tea.Quit
		}
	case error:
		m.done = true
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m ProgressModel) View() string {
	const total = 30
	filled := int(m.progress * total)
	bar := ""
	for i := 0; i < total; i++ {
		if i < filled {
			bar += lipgloss.NewStyle().Background(Green).Render(" ")
		} else if i == filled {
			bar += lipgloss.NewStyle().Background(Yellow).Render(">")
		} else {
			bar += lipgloss.NewStyle().Background(lipgloss.Color("#333")).Render(" ")
		}
	}
	label := fmt.Sprintf("%3.0f%%", m.progress*100)
	return fmt.Sprintf("%s %s %s", m.title, bar, label)
}
