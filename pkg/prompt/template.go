/*
Package prompt handles the generation and templating of instructions for AI agents.

It uses standard Go templates to inject repository-specific context (like the
repository name and branch) into a set of behavioral mandates, ensuring the
AI agent understands its constraints and goals within the multi-gitter pipeline.
*/
package prompt

import (
	"bytes"
	"fmt"
	"text/template"
)

// DefaultAgentTemplate is the markdown template used to instruct the AI agent.
const DefaultAgentTemplate = `
# multi-gitter-agent: AI Orchestration Task

## 🛠 Context
You are an AI agent acting as a specialized "worker" in a multi-repository refactoring pipeline.
- **Orchestrator**: multi-gitter
- **Target Repository**: {{ .Repository }}
{{- if .BaseBranch }}
- **Base Branch**: {{ .BaseBranch }}
{{- end }}
- **Operation Mode**: {{ if .DryRun }}DRY RUN (Simulated){{ else }}LIVE (Changes will be pushed){{ end }}

## 🎯 Objective
{{ .UserPrompt }}

## 🛡 Security & Safety Mandates
1. **No Secrets**: NEVER hardcode API keys, tokens, or credentials. Do not log sensitive data.
2. **Scope**: Only modify files directly related to the objective. Avoid unrelated "cleanups" unless specifically requested.
3. **Validation**: Ensure that any code you write is syntactically correct and follows the repository's established style.
4. **Environment**: You are running in a temporary clone of the repository. Your changes will be evaluated for a Pull Request.

## ⚙️ Operational Workflow
1. **Discovery**: Use your tools (grep, ls, find) to determine if this repository actually needs the requested change.
2. **Actionability**: If the objective is NOT applicable to this repository, **do not modify any files**.
3. **Triggering**: multi-gitter will only create a Pull Request if it detects file modifications. No modification = No PR.
4. **Lifecycle**: Your work is automatically committed and pushed to a new branch. You do NOT need to run git commit, git push, or manage PRs yourself.

## 📝 Execution Instructions
- Perform a thorough analysis of the codebase.
- Apply the changes surgically.
- If you encounter an error or ambiguity, prioritize safety and avoid making partial or broken changes.
`

// Context holds the information needed to render the agent prompt.
type Context struct {
	UserPrompt string
	Repository string
	BaseBranch string
	DryRun     bool
}

// Registry manages available prompt templates.
type Registry struct {
	templates map[string]string
}

// NewRegistry creates a new template registry with the default template.
func NewRegistry() *Registry {
	r := &Registry{
		templates: make(map[string]string),
	}
	r.Register("default", DefaultAgentTemplate)
	return r
}

// Register adds a new template to the registry.
func (r *Registry) Register(name, content string) {
	r.templates[name] = content
}

// Render generates the final prompt string using the named template and context.
func (r *Registry) Render(name string, ctx Context) (string, error) {
	content, ok := r.templates[name]
	if !ok {
		// If not in registry, assume name is a raw template or use default
		content = r.templates["default"]
	}

	tmpl, err := template.New("prompt").Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
