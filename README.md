# multi-gitter-agent 🤖

A high-performance CLI tool that combines [multi-gitter](https://github.com/lindell/multi-gitter) with AI agents (Gemini, Claude Code, Copilot) to perform natural language refactors across one or many repositories.

## Overview

`multi-gitter-agent` allows you to orchestrate complex code modifications across your entire organization using AI. It handles the cloning, branching, and Pull Request creation while delegating the logic to your favorite AI agent.

## Installation

### Shell Installer (Linux & macOS)
Install the latest binary automatically:
```bash
curl -sSL https://raw.githubusercontent.com/juftin/multi-gitter-agent/main/install.sh | bash
```

### Via Go
```bash
go install github.com/juftin/multi-gitter-agent@latest
```

## CLI Flags

- `-p, --prompt`: The natural language instruction for the agent.
- `--prompt-file`: Path to a file containing the instruction.
- `-a, --agent`: (**Required**) The AI agent to use (`gemini`, `claude`, `copilot`). Can also be set via `MULTI_GITTER_AGENT` environment variable.
- `-y, --yolo`: Automatically approve all AI actions.
- `-i, --interactive`: Take manual decisions before committing changes.
- `-C, --concurrent`: Number of repositories to process simultaneously (Interactive mode requires `-C 1`).
- `--org, --repo, --user, --topic`: Standard `multi-gitter` repository discovery flags.
- `--repo-search, --code-search`: GitHub search discovery flags.
- `--repo-include, --repo-exclude`: Regex-based repository filtering.
- `--draft`: Create Pull Requests as drafts.
- `-t, --token`: Your platform access token (e.g., `GITHUB_TOKEN`).

## Usage Examples

### Bulk Dependency Update (Automated)
Update a specific dependency across all repositories in an organization using Gemini:
```bash
multi-gitter-agent run --agent gemini \
  --prompt "Update github.com/stretchr/testify to v1.10.0 in go.mod" \
  --org my-org --yolo
```

### Targeted Refactor via Code Search (Interactive)
Use GitHub Code Search to find repositories containing a specific pattern and refactor them with Claude:
```bash
multi-gitter-agent run --agent claude \
  --code-search "filename:package.json lodash" \
  --prompt "Replace lodash with native Array methods where possible" \
  --interactive
```

### Advanced Repository Discovery
Find repositories by GitHub Topics and exclude legacy projects using Regex:
```bash
multi-gitter-agent run --agent gemini \
  --topic production --topic microservice \
  --repo-exclude "^legacy-.*" \
  --prompt "Add a HEALTHCHECK instruction to the Dockerfile" \
  --yolo
```

### Targeted Repository Search
Use GitHub's repository search syntax to target specific languages or metadata:
```bash
multi-gitter-agent run --agent gemini \
  --repo-search "topic:react language:typescript" \
  --prompt "Migrate from TSLint to ESLint" \
  --draft --yolo
```

### Organization-wide Security Audit
Run a security audit across all repositories and only create PRs for actionable fixes:
```bash
multi-gitter-agent run --agent gemini \
  --prompt "Check for hardcoded API keys and move them to environment variables" \
  --org my-company --concurrent 5 --yolo
```

### Complex Refactor with Prompt File
Provide detailed, multi-step instructions from a Markdown file:
```bash
multi-gitter-agent run --agent copilot \
  --prompt-file ./docs/refactor-plan.md \
  --user my-github-username --branch feature/api-v2-migration
```

## Features

- **Integrated Orchestration**: Natively wraps `multi-gitter` for seamless multi-repo management.
- **Agent Support**: Built-in support for **Gemini CLI**, **Claude Code**, and **GitHub Copilot**.
- **Interactive TUI**: A polished Bubble Tea interface for managing per-repository tasks.
- **YOLO Mode**: Fully automated execution for bulk refactoring.
- **Robust TTY Handling**: Custom PTY wrapper ensures stability for complex AI CLIs.
- **Context-Rich Prompts**: Automatically provides agents with repository metadata and safety mandates.

## How it works

When you execute `multi-gitter-agent run`, the following orchestration occurs:

1. **Repository Discovery**: The tool uses the provided flags (`--org`, `--user`, `--repo`, or `--code-search`) to identify all target repositories.
2. **Cloning**: `multi-gitter-agent` (via the embedded `multi-gitter` library) creates temporary clones of each repository and checks out a new feature branch.
3. **Agent Invocation**: For each repository, the tool launches an **Internal Agent Runner**. This runner:
   - Sets up a stable terminal environment using a **Pseudo-Terminal (PTY)**.
   - Renders a rich Markdown prompt containing your instructions and the repository's context.
   - Executes your chosen AI agent (Gemini, Claude, or Copilot) within that PTY.
4. **AI Dialogue**: The AI agent analyzes the code and applies modifications. In `--interactive` mode, you can review and respond to the agent's thoughts. In `--yolo` mode, all agent actions are automatically approved.
5. **Change Detection**: After the agent finishes, the tool checks for file modifications.
6. **PR Creation**: If changes were made, the tool automatically commits the work, pushes the branch, and opens a Pull Request on your target platform (GitHub, GitLab, etc.).

## Safety

- **Dry Run**: Use `--dry-run` to see what changes the AI would make without pushing anything.
- **Security Mandates**: The tool explicitly instructs agents to never hardcode secrets or modify files outside the requested scope.
