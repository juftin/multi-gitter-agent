package agent

import (
	"context"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/juftin/multi-gitter-agent/pkg/prompt"
)

// Options holds the configuration for running an AI agent.
type Options struct {
	Repo         string
	BaseBranch   string
	DryRun       bool
	UserPrompt   string
	AgentCommand string
	AgentArgs    []string
	Interactive  bool
	Yolo         bool
	Silent       bool
}

// Run executes the AI agent for a given repository.
func Run(ctx context.Context, opts Options) error {
	if opts.Silent {
		return runSilent(ctx, opts)
	}

	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	var tty *os.File
	if err == nil {
		tty = f
		defer tty.Close()
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m := model{
		opts:    opts,
		spinner: s,
		tty:     tty,
	}

	var p *tea.Program
	if tty != nil {
		p = tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	} else {
		p = tea.NewProgram(m)
	}

	_, err = p.Run()
	return err
}

func runSilent(ctx context.Context, opts Options) error {
	promptRegistry := prompt.NewRegistry()
	promptContent, err := promptRegistry.Render("default", prompt.Context{
		UserPrompt: opts.UserPrompt,
		Repository: opts.Repo,
		BaseBranch: opts.BaseBranch,
		DryRun:     opts.DryRun,
	})
	if err != nil {
		return err
	}

	registry := NewRegistry()
	provider, err := registry.Get(opts.AgentCommand)
	if err != nil {
		provider = &GenericAgent{Command: opts.AgentCommand}
	}

	llmCmd, err := provider.BuildCommand(ctx, promptContent, opts)
	if err != nil {
		return err
	}
	PrepareEnv(llmCmd)

	// Suppress output in silent mode
	llmCmd.Stdout = nil
	llmCmd.Stderr = nil

	return llmCmd.Run()
}

type agentFinishedMsg struct{ err error }

func (m model) runAgent() tea.Cmd {
	promptRegistry := prompt.NewRegistry()
	promptContent, err := promptRegistry.Render("default", prompt.Context{
		UserPrompt: m.opts.UserPrompt,
		Repository: m.opts.Repo,
		BaseBranch: m.opts.BaseBranch,
		DryRun:     m.opts.DryRun,
	})
	if err != nil {
		return func() tea.Msg {
			return agentFinishedMsg{err}
		}
	}

	registry := NewRegistry()
	provider, err := registry.Get(m.opts.AgentCommand)
	if err != nil {
		provider = &GenericAgent{Command: m.opts.AgentCommand}
	}

	// We use background context for the command since it's an interactive TUI step
	llmCmd, err := provider.BuildCommand(context.Background(), promptContent, m.opts)
	if err != nil {
		return func() tea.Msg {
			return agentFinishedMsg{err}
		}
	}
	PrepareEnv(llmCmd)

	ptyCmd := &ptyCommand{
		Cmd: llmCmd,
		tty: m.tty,
	}

	// THE FIX: Use tea.Exec to properly suspend the TUI and let the PTY run.
	return tea.Exec(ptyCmd, func(err error) tea.Msg {
		return agentFinishedMsg{err}
	})
}
