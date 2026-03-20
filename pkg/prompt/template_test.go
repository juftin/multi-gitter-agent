package prompt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry_Render(t *testing.T) {
	r := NewRegistry()

	ctx := Context{
		UserPrompt: "Update documentation",
		Repository: "my-org/my-repo",
		BaseBranch: "main",
		DryRun:     true,
	}

	t.Run("render default template", func(t *testing.T) {
		got, err := r.Render("default", ctx)
		require.NoError(t, err)
		assert.Contains(t, got, "# multi-gitter-agent: AI Task")
		assert.Contains(t, got, "Update documentation")
		assert.Contains(t, got, "my-org/my-repo")
		assert.Contains(t, got, "DRY RUN")
	})

	t.Run("render with live mode", func(t *testing.T) {
		ctxLive := ctx
		ctxLive.DryRun = false
		got, err := r.Render("default", ctxLive)
		require.NoError(t, err)
		assert.Contains(t, got, "LIVE (Changes will be pushed)")
	})

	t.Run("render custom template", func(t *testing.T) {
		custom := "Custom: {{ .UserPrompt }} on {{ .Repository }}"
		r.Register("custom", custom)
		got, err := r.Render("custom", ctx)
		require.NoError(t, err)
		assert.Equal(t, "Custom: Update documentation on my-org/my-repo", got)
	})

	t.Run("fallback to default for unknown template", func(t *testing.T) {
		got, err := r.Render("non-existent", ctx)
		require.NoError(t, err)
		assert.Contains(t, got, "# multi-gitter-agent: AI Task")
	})

	t.Run("invalid template error", func(t *testing.T) {
		r.Register("invalid", "Invalid {{ .NonExistent }}")
		_, err := r.Render("invalid", ctx)
		assert.Error(t, err)
	})
}

func TestRenderPRBody(t *testing.T) {
	ctx := Context{
		UserPrompt: "Add tests",
	}
	body := RenderPRBody(ctx, "gemini", "Full prompt text")

	assert.Contains(t, body, "## 🤖 multi-gitter-agent PR")
	assert.Contains(t, body, "using the **gemini** agent")
	assert.Contains(t, body, "Add tests")
	assert.Contains(t, body, "Full prompt text")
	assert.Contains(t, body, "<details>")
}

func TestRenderTemplateFile(t *testing.T) {
	t.Run("renders template file", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "prompt-template-*.tmpl")
		require.NoError(t, err)
		defer f.Close()

		_, err = f.WriteString("Task: {{ .UserPrompt }} in {{ .Repository }}")
		require.NoError(t, err)

		got, err := RenderTemplateFile(f.Name(), Context{
			UserPrompt: "Do thing",
			Repository: "org/repo",
		})
		require.NoError(t, err)
		assert.Equal(t, "Task: Do thing in org/repo", got)
	})

	t.Run("returns error on missing file", func(t *testing.T) {
		_, err := RenderTemplateFile("/does/not/exist.tmpl", Context{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read template file")
	})
}
