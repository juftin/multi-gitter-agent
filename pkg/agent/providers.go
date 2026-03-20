package agent

import (
	"context"
	"os"
	"os/exec"
)

// GeminiAgent provides support for the Gemini CLI.
type GeminiAgent struct{}

func (a *GeminiAgent) ID() string          { return "gemini" }
func (a *GeminiAgent) DisplayName() string { return "Gemini CLI" }
func (a *GeminiAgent) Description() string { return "Google's AI model via the Gemini CLI" }

func (a *GeminiAgent) BuildCommand(ctx context.Context, prompt string, opts Options) (*exec.Cmd, error) {
	var args []string
	if opts.Interactive {
		args = append(args, "--prompt-interactive", prompt)
	} else {
		args = append(args, "--prompt", prompt, "--output-format", "text")
	}
	if opts.Yolo {
		args = append(args, "--yolo", "--accept-raw-output-risk")
	}
	args = append(args, opts.AgentArgs...)
	return exec.CommandContext(ctx, "gemini", args...), nil
}

// ClaudeAgent provides support for Anthropic's Claude Code.
type ClaudeAgent struct{}

func (a *ClaudeAgent) ID() string          { return "claude" }
func (a *ClaudeAgent) DisplayName() string { return "Claude Code" }
func (a *ClaudeAgent) Description() string { return "Anthropic's Claude model for code editing" }

func (a *ClaudeAgent) BuildCommand(ctx context.Context, prompt string, opts Options) (*exec.Cmd, error) {
	var args []string
	if !opts.Interactive {
		args = append(args, "-p", prompt, "--output-format", "text")
	} else {
		args = append(args, prompt)
	}
	if opts.Yolo {
		args = append(args, "--allow-dangerously-skip-permissions", "--dangerously-skip-permissions")
	}
	args = append(args, opts.AgentArgs...)
	return exec.CommandContext(ctx, "claude", args...), nil
}

// CopilotAgent provides support for GitHub Copilot CLI.
type CopilotAgent struct{}

func (a *CopilotAgent) ID() string          { return "copilot" }
func (a *CopilotAgent) DisplayName() string { return "GitHub Copilot" }
func (a *CopilotAgent) Description() string {
	return "AI programming assistant via the gh copilot extension"
}

func (a *CopilotAgent) BuildCommand(ctx context.Context, prompt string, opts Options) (*exec.Cmd, error) {
	binary := "copilot"
	var args []string
	if opts.AgentCommand == "gh copilot" {
		binary = "gh"
		args = append(args, "copilot")
	}
	args = append(args, "-p", prompt)
	if opts.Yolo {
		args = append(args, "--yolo")
	}
	args = append(args, opts.AgentArgs...)
	return exec.CommandContext(ctx, binary, args...), nil
}

// GenericAgent is used when an unknown command is passed as an agent.
type GenericAgent struct {
	Command string
}

func (a *GenericAgent) ID() string          { return a.Command }
func (a *GenericAgent) DisplayName() string { return a.Command }
func (a *GenericAgent) Description() string { return "Custom agent command" }

func (a *GenericAgent) BuildCommand(ctx context.Context, prompt string, opts Options) (*exec.Cmd, error) {
	args := []string{prompt}
	args = append(args, opts.AgentArgs...)
	return exec.CommandContext(ctx, a.Command, args...), nil
}

// PrepareEnv adds common terminal-related environment variables.
func PrepareEnv(cmd *exec.Cmd) {
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"FORCE_COLOR=1",
		"CLAUDE_CODE_OUTPUT=text",
		"GEMINI_OUTPUT=text",
	)
}
