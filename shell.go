// Package shellout provides a thin wrapper around os/exec to improve testability and ergonomics
// when executing simple shell commands programmatically.
//
// Note that the improved ergonomics come with reduced flexibility. This package is not appropriate
// for running processes that require dynamic interactions with stdin and stdout or commands that
// produce large volumes of output.
package shellout

import (
	"bytes"
	"errors"
	"io"
	"os/exec"

	"github.com/ttd2089/tyers"
)

// ErrCommandNotFound is returned when a requested command can not be resolved from the PATH.
var ErrCommandNotFound error = errors.New("ErrCommandNotFound")

// ErrCommandProcessFailed is return when there was a problem executing the command process. Note
// that unlike with the os/exec package, a process exiting with a non-zero status code will not be
// represented as an error.
var ErrCommandProcessFailed error = errors.New("ErrCommandProcessFailed")

// A Cmd contains the information required to start a command process.
//
// Except where documentation states otherwise the semantics of each property of a Cmd are defined
// by the semantics of the corrsponding property on the exec.Cmd type.
//
// Refer to exec.Command for the semantics of the Command and Args properties.
type Cmd struct {

	// Command is the name or path of the command to be executed.
	Command string

	// Args are the arguments to pass to the process.
	Args []string

	// Env specifies the environment of the process.
	Env []string

	// Dir specifies the directory to run the command in.
	Dir string

	// Stdin specifies the process's standard input.
	Stdin io.Reader
}

// A Result represents the outcome of running a command.
type Result struct {

	// The exit code from the command that was run.
	ExitCode int

	// Stdout contains the data written to stdout by the command process.
	Stdout *bytes.Buffer

	// Stdout contains the data written to stderr by the command process.
	Stderr *bytes.Buffer
}

// Run executes the Run method on a default instance of Shell.
func Run(cmd Cmd) (Result, error) {
	return defaultShell.Run(cmd)
}

// A Shell provides a simple CLI-like interface for executing processes.
type Shell interface {

	// Run executes the specified Cmd and captures the output as a Result.
	//
	// If the process cannot be started then an error is returned; otherwise the Result is
	// populated with the results. Note that unlike the with os/exec package, a process exiting
	// with a non-zero status code will not be represented as an error.
	Run(cmd Cmd) (Result, error)
}

// New returns an instance of Shell.
func New() Shell {
	return defaultShell
}

type shell struct{}

var defaultShell shell

func (_ shell) Run(cmd Cmd) (Result, error) {
	proc := exec.Command(cmd.Command, cmd.Args...)
	if proc.Err != nil {
		return Result{}, tyers.As(ErrCommandNotFound, proc.Err)
	}
	result := Result{
		Stdout: new(bytes.Buffer),
		Stderr: new(bytes.Buffer),
	}
	proc.Env = cmd.Env
	proc.Dir = cmd.Dir
	proc.Stdin = cmd.Stdin
	proc.Stdout = result.Stdout
	proc.Stderr = result.Stderr
	err := proc.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		return Result{}, tyers.As(ErrCommandProcessFailed, err)
	}
	return result, nil
}
