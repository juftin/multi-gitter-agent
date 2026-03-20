package agent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	t.Run("default agents registered", func(t *testing.T) {
		agents := r.List()
		assert.Len(t, agents, 3)

		ids := []string{}
		for _, a := range agents {
			ids = append(ids, a.ID())
		}
		assert.Contains(t, ids, "gemini")
		assert.Contains(t, ids, "claude")
		assert.Contains(t, ids, "copilot")
	})

	t.Run("get agent", func(t *testing.T) {
		a, err := r.Get("gemini")
		require.NoError(t, err)
		assert.Equal(t, "gemini", a.ID())
	})

	t.Run("register custom agent", func(t *testing.T) {
		r.Register(&GenericAgent{Command: "my-ai"})
		a, err := r.Get("my-ai")
		require.NoError(t, err)
		assert.Equal(t, "my-ai", a.ID())
	})
}

func TestAgent_BuildCommand(t *testing.T) {
	ctx := context.Background()
	prompt := "Say hello"
	opts := Options{
		Yolo:        true,
		Interactive: false,
	}

	t.Run("gemini yolo", func(t *testing.T) {
		a := &GeminiAgent{}
		cmd, err := a.BuildCommand(ctx, prompt, opts)
		require.NoError(t, err)
		assert.Contains(t, cmd.Args, "--yolo")
		assert.Contains(t, cmd.Args, "--accept-raw-output-risk")
		assert.Contains(t, cmd.Args, prompt)
	})

	t.Run("claude non-interactive", func(t *testing.T) {
		a := &ClaudeAgent{}
		cmd, err := a.BuildCommand(ctx, prompt, opts)
		require.NoError(t, err)
		assert.Contains(t, cmd.Args, "-p")
		assert.Contains(t, cmd.Args, "--output-format")
		assert.Contains(t, cmd.Args, "text")
	})

	t.Run("claude yolo", func(t *testing.T) {
		a := &ClaudeAgent{}
		cmd, err := a.BuildCommand(ctx, prompt, opts)
		require.NoError(t, err)
		assert.Contains(t, cmd.Args, "--allow-dangerously-skip-permissions")
		assert.Contains(t, cmd.Args, "--dangerously-skip-permissions")
	})

	t.Run("copilot yolo", func(t *testing.T) {
		a := &CopilotAgent{}
		cmd, err := a.BuildCommand(ctx, prompt, opts)
		require.NoError(t, err)
		assert.Contains(t, cmd.Args, "--yolo")
		assert.Contains(t, cmd.Args, "-p")
	})
}
