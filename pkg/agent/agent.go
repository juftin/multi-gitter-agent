/*
Package agent provides the core orchestration logic for executing AI assistants
(like Gemini, Claude, and Copilot) within a terminal environment.

It defines an extensible Agent interface and a Registry for managing them. It also
includes the necessary TUI (Terminal User Interface) and PTY (Pseudo-Terminal)
wrappers to ensure stable, interactive execution of child processes.
*/
package agent

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
)

// Agent defines the behavior and metadata for an AI provider.
type Agent interface {
	// ID returns the unique identifier for the agent (e.g., "gemini").
	ID() string
	// DisplayName returns a human-readable name for the agent.
	DisplayName() string
	// Description returns a brief explanation of the agent's purpose.
	Description() string
	// BuildCommand constructs the execution command for this agent.
	BuildCommand(ctx context.Context, prompt string, opts Options) (*exec.Cmd, error)
}

// Registry manages the collection of available AI agents.
type Registry struct {
	agents map[string]Agent
}

// NewRegistry creates a new agent registry and registers default providers.
func NewRegistry() *Registry {
	r := &Registry{
		agents: make(map[string]Agent),
	}
	r.Register(&GeminiAgent{})
	r.Register(&ClaudeAgent{})
	r.Register(&CopilotAgent{})
	return r
}

// Register adds a new agent to the registry.
func (r *Registry) Register(a Agent) {
	r.agents[a.ID()] = a
}

// Get retrieves an agent by its ID.
func (r *Registry) Get(id string) (Agent, error) {
	a, ok := r.agents[id]
	if !ok {
		return nil, fmt.Errorf("unknown agent: %s", id)
	}
	return a, nil
}

// List returns all registered agents sorted by their ID.
func (r *Registry) List() []Agent {
	var agents []Agent
	for _, a := range r.agents {
		agents = append(agents, a)
	}
	sort.Slice(agents, func(i, j int) bool {
		return agents[i].ID() < agents[j].ID()
	})
	return agents
}
