package multigitter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	opts := Options{
		AgentScript: "test-agent",
		Orgs:        []string{"my-org"},
		PRTitle:     "Test PR",
		DryRun:      true,
	}

	assert.Equal(t, "test-agent", opts.AgentScript)
	assert.Equal(t, []string{"my-org"}, opts.Orgs)
	assert.True(t, opts.DryRun)
}

func TestRun_Validation(t *testing.T) {
	// Running without a PR title should fail
	opts := Options{
		AgentScript: "echo hello",
		Repos:       []string{"juftin/multi-gitter-agent"},
	}

	err := Run(context.Background(), opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pull request title or commit message must be set")
}
