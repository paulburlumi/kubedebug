package kubedebug

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

// Command executor
type Command struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewCommand constructs a [Command]
func NewCommand(stdin io.Reader, stdout, stderr io.Writer) *Command {
	return &Command{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Output runs the command and returns its standard output.
func (c *Command) Output(name string, arg ...string) (string, error) {
	cmd := exec.CommandContext(context.Background(), name, arg...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%v failed: %w", cmd.Args, err)
	}
	return string(out), nil
}

// Run starts the specified command and waits for it to complete.
func (c *Command) Run(redirect bool, name string, arg ...string) error {
	cmd := exec.CommandContext(context.Background(), name, arg...)
	if redirect {
		cmd.Stdin = c.stdin
		cmd.Stdout = c.stdout
		cmd.Stderr = c.stderr
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v failed: %w", cmd.Args, err)
	}
	return nil
}
