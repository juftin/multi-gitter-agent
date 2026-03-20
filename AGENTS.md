# AGENTS.md

This file contains foundational mandates and instructions for AI agents working on `multi-gitter-agent`.

## Core Architecture

The application is structured to decouple the CLI logic, the multi-gitter orchestration, and the AI agent execution.

- **`cmd/multi-gitter-agent/main.go`**: The entry point. It defines the Cobra CLI commands (`run` and the hidden `agent` subcommand).
- **`pkg/agent/`**: The core execution engine for AI models.
  - `agent.go`: Defines the extensible `Agent` interface and the `Registry` pattern.
  - `providers.go`: Implementations for specific AI CLIs (Gemini, Claude, Copilot).
  - `runner.go`: Orchestration logic for building and executing agent commands, including TUI and Silent modes.
  - `tui.go`: The Bubble Tea interactive terminal dashboard.
  - `pty.go`: A Pseudo-Terminal (PTY) wrapper essential for terminal stability (prevents `kqueue`/`isTTY` crashes in runtimes like Bun).
- **`pkg/multigitter/`**: Library-level wrapper embedding `lindell/multi-gitter`. It maps our CLI flags directly to the underlying engine.
- **`pkg/prompt/`**: Markdown prompt templates and Go templating logic via a `Registry`.

## Engineering Mandates

1. **Terminal Stability**: Always use the PTY wrapper (`pkg/agent/pty.go`) when launching interactive agents. This satisfies strict `isTTY` checks required by modern JS runtimes.
2. **Concurrency Safety**: If `concurrent > 1`, agents must run via `runSilent()` (no TUI, suppressed `stdout`/`stderr`) to prevent terminal display corruption.
3. **Agent Mappings & YOLO**:
   - **Gemini**: Use `--prompt` (automated) or `--prompt-interactive` (interactive). Use `--yolo` and `--accept-raw-output-risk` for auto-approval.
   - **Claude**: Use `-p` (prompt). Use `--allow-dangerously-skip-permissions` and `--dangerously-skip-permissions` for auto-approval.
   - **Copilot**: Use `-p` (prompt) and `--yolo` for auto-approval.
4. **Actionability**: The prompt template (`pkg/prompt/template.go`) explicitly instructs workers to only modify files if the task is applicable. No modifications = No Pull Request.
5. **Self-Contained Executable**: Do not rely on external `multi-gitter` binaries. Always execute via `cmd.RootCmd()` embedded from the library.

## Operational Flow

1. User executes `multi-gitter-agent run`.
2. The `multigitter` package translates options and invokes the internal `multi-gitter` library.
3. The library discovers repositories and, for each one, invokes the `multi-gitter-agent agent` subcommand as a child process.
4. The `agent` subcommand hijacks the terminal (using `/dev/tty`), initializes the Bubble Tea TUI, and launches the AI agent inside a PTY.
5. If the agent modifies files, the parent `multi-gitter` process detects the diff, commits the changes, and opens a Pull Request.

## Development Standards

- Use `gofmt` and `goimports`.
- Document all exported symbols, structs, and packages.
- Prioritize modularity in `pkg/agent/` to allow for adding new LLM providers effortlessly.
- Use `testify` for all unit testing.
