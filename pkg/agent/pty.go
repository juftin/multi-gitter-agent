package agent

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// ptyCommand wraps an exec.Cmd to be run within a pseudo-terminal (PTY).
// This is necessary for some AI agents (like Claude Code/Bun) that require
// a true TTY to function correctly without crashing (e.g., kqueue errors).
type ptyCommand struct {
	*exec.Cmd
	tty *os.File
}

// SetStdin is a no-op as the PTY handles standard input.
func (c *ptyCommand) SetStdin(r io.Reader) {}

// SetStdout is a no-op as the PTY handles standard output.
func (c *ptyCommand) SetStdout(w io.Writer) {}

// SetStderr is a no-op as the PTY handles standard error.
func (c *ptyCommand) SetStderr(w io.Writer) {}

// Run starts the command within a PTY, copies I/O between the real TTY
// and the PTY, and waits for the command to complete.
func (c *ptyCommand) Run() error {
	ptmx, err := pty.Start(c.Cmd)
	if err != nil {
		return err
	}
	defer ptmx.Close()

	if c.tty != nil {
		// Inherit size from the real terminal
		if err := pty.InheritSize(c.tty, ptmx); err != nil {
			fmt.Printf("Error inheriting terminal size: %v\n", err)
		}
	}

	var oldState *term.State
	if c.tty != nil {
		// Put the real TTY into raw mode so keystrokes pass through to the PTY
		oldState, err = term.MakeRaw(int(c.tty.Fd()))
		if err == nil {
			defer func() {
				_ = term.Restore(int(c.tty.Fd()), oldState)
			}()
		}

		// Copy input from the real TTY to the PTY
		go func() {
			_, _ = io.Copy(ptmx, c.tty)
		}()
		// Copy output from the PTY to the real TTY
		go func() {
			_, _ = io.Copy(c.tty, ptmx)
		}()
	} else {
		// Fallback if no TTY is provided
		go func() {
			_, _ = io.Copy(ptmx, os.Stdin)
		}()
		go func() {
			_, _ = io.Copy(os.Stdout, ptmx)
		}()
	}

	return c.Cmd.Wait()
}
