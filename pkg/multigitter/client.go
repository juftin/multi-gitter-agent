/*
Package multigitter provides a native wrapper around the github.com/lindell/multi-gitter library.

By embedding the multi-gitter RootCmd directly, this package allows the
multi-gitter-agent to orchestrate multi-repository cloning, branching, and PR
creation without requiring the user to install the multi-gitter CLI separately.
*/
package multigitter

import (
	"context"
	"fmt"
	"os"

	"github.com/lindell/multi-gitter/cmd"
)

// Options holds the configuration for executing the multi-gitter logic.
type Options struct {
	AgentScript      string
	Orgs             []string
	Repos            []string
	Users            []string
	Topics           []string
	RepoSearch       string
	CodeSearch       string
	RepoInclude      string
	RepoExclude      string
	SkipForks        bool
	Fork             bool
	Platforms        string
	Token            string
	Branch           string
	BaseBranch       string
	ConflictStrategy string
	Draft            bool
	PRTitle          string
	PRBody           string
	DryRun           bool
	Interactive      bool
	Concurrent       int
}

// Run executes the multi-gitter root command directly from the library.
func Run(ctx context.Context, opts Options) error {
	root := cmd.RootCmd()

	// Prepare arguments for multi-gitter
	args := []string{"run", opts.AgentScript}

	if len(opts.Orgs) > 0 {
		for _, org := range opts.Orgs {
			args = append(args, "--org", org)
		}
	}
	if len(opts.Repos) > 0 {
		for _, repo := range opts.Repos {
			args = append(args, "--repo", repo)
		}
	}
	if len(opts.Users) > 0 {
		for _, user := range opts.Users {
			args = append(args, "--user", user)
		}
	}
	if len(opts.Topics) > 0 {
		for _, topic := range opts.Topics {
			args = append(args, "--topic", topic)
		}
	}
	if opts.RepoSearch != "" {
		args = append(args, "--repo-search", opts.RepoSearch)
	}
	if opts.CodeSearch != "" {
		args = append(args, "--code-search", opts.CodeSearch)
	}
	if opts.RepoInclude != "" {
		args = append(args, "--repo-include", opts.RepoInclude)
	}
	if opts.RepoExclude != "" {
		args = append(args, "--repo-exclude", opts.RepoExclude)
	}
	if opts.SkipForks {
		args = append(args, "--skip-forks")
	}
	if opts.Fork {
		args = append(args, "--fork")
	}
	if opts.Platforms != "" {
		args = append(args, "--platform", opts.Platforms)
	}
	if opts.Token != "" {
		args = append(args, "--token", opts.Token)
	}
	if opts.Branch != "" {
		args = append(args, "--branch", opts.Branch)
	}
	if opts.BaseBranch != "" {
		args = append(args, "--base-branch", opts.BaseBranch)
	}
	if opts.ConflictStrategy != "" {
		args = append(args, "--conflict-strategy", opts.ConflictStrategy)
	}
	if opts.Draft {
		args = append(args, "--draft")
	}
	if opts.PRTitle != "" {
		args = append(args, "--pr-title", opts.PRTitle)
	}
	if opts.PRBody != "" {
		args = append(args, "--pr-body", opts.PRBody)
	}
	if opts.DryRun {
		args = append(args, "--dry-run")
	}
	if opts.Interactive {
		args = append(args, "--interactive")
	}
	if opts.Concurrent > 0 {
		args = append(args, "--concurrent", fmt.Sprintf("%d", opts.Concurrent))
	}

	// Set multi-gitter's internal arguments
	root.SetArgs(args)

	// Since we are running in-process, multi-gitter will use os.Stdout/Err/In by default.
	// We ensure they are wired up correctly.
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)
	root.SetIn(os.Stdin)

	return root.ExecuteContext(ctx)
}
