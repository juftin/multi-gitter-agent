# multi-gitter-agent

A CLI tool that combines the power of [multi-gitter](https://github.com/lindell/multi-gitter) with LLM agents (like Claude Code/Gemini CLI) to run natural language prompts across one or many repositories.

## Overview

`multi-gitter-agent` allows you to perform complex refactors or audits across multiple repositories using AI agents. It handles the repository cloning, branching, and PR creation (via `multi-gitter`) while delegating the code modification logic to an AI agent.

## Features

- **Bulk AI Refactoring**: Run an LLM prompt against multiple repos simultaneously.
- **Repository Management**: Leverages `multi-gitter` for robust git operations.
- **Agent Integration**: Seamlessly passes context to AI agents for precise modifications.

## Installation

```bash
go install github.com/juftin/multi-gitter-agent@latest
```

## Usage

```bash
multi-gitter-agent run "Replace all instances of DeprecatedMethod() with NewMethod()" --org my-org
```
