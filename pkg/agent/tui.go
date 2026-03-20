package agent

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)
	repoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true)
	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

// model represents the state of the Bubble Tea TUI.
type model struct {
	opts     Options
	spinner  spinner.Model
	quitting bool
	err      error
	done     bool
	tty      *os.File
	width    int
	height   int
}

// Init initializes the Bubble Tea model.
func (m model) Init() tea.Cmd {
	if m.opts.Yolo {
		return tea.Batch(m.spinner.Tick, m.runAgent())
	}
	return m.spinner.Tick
}

// Update handles incoming messages and updates the model's state.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if !m.done {
				return m, m.runAgent()
			}
		}
	case agentFinishedMsg:
		m.done = true
		m.quitting = true // Also set quitting to true to ensure empty final render
		if msg.err != nil {
			m.err = msg.err
		}
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

// View renders the current state of the model as a string.
func (m model) View() string {
	if m.quitting || m.done {
		return ""
	}
	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n", m.err)
	}

	style := lipgloss.NewStyle().Width(m.width)

	var s string
	s += titleStyle.Render("multi-gitter-agent") + "\n\n"
	s += style.Render(fmt.Sprintf("Repository: %s", repoStyle.Render(m.opts.Repo))) + "\n"
	s += style.Render(fmt.Sprintf("Prompt:     %s", m.opts.UserPrompt)) + "\n\n"

	s += m.spinner.View() + statusStyle.Render(" Press Enter to start agent or 'q' to skip repo...")

	return s
}
