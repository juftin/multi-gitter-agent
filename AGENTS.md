# AGENTS.md

This file contains foundational mandates and instructions for AI agents (Gemini CLI, Claude Code, etc.) working on this project. These instructions take precedence over general workflows.

## Core Mandates

- **Language**: Use Go (Golang) for all core logic.
- **CLI Framework**: Use `spf13/cobra` for command management and `spf13/viper` for configuration.
- **Error Handling**: Use structured error handling. Wrap errors with context where appropriate.
- **Testing**: Every feature MUST have corresponding unit tests. Use `testify` for assertions.
- **Integration**: Leverage `multi-gitter` functionality where possible.

## Project Vision

`multi-gitter-agent` is designed to be a bridge between multi-repository management and LLM-driven code modification.

1. **Repository Discovery**: Identify repositories to act upon (via `multi-gitter` patterns).
2. **Context Gathering**: Prepare the repository state for the agent.
3. **Agent Execution**: Invoke an LLM agent with the user's prompt.
4. **Lifecycle Management**: Commit, branch, and PR management using `multi-gitter`.

## Architecture

- `cmd/`: CLI entry points and command definitions.
- `pkg/agent/`: Logic for interacting with LLM agents.
- `pkg/multigitter/`: Wrapper/integration logic for `multi-gitter`.
- `pkg/prompt/`: Management and templates for natural language prompts.

## Development Standards

- Follow standard Go project layout (`cmd/`, `pkg/`, `internal/`).
- Use `gofmt` and `goimports` for code formatting.
- Document all public functions and types.
- Prioritize modularity to allow for different LLM backends (OpenAI, Anthropic, Gemini).

## Success Criteria for Tasks

- Code is idiomatic Go.
- Command-line interface is intuitive and follows POSIX standards.
- Tests cover both success and failure paths.
- Performance: Minimize overhead when orchestrating across many repositories.
