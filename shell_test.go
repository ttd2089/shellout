package shellout

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {

	t.Run("Non-existent command returns ErrCommandNotFound", func(t *testing.T) {
		nonExistentCommand := "icantbelievethisisacommandinyourenvironment"
		if _, err := exec.LookPath(nonExistentCommand); err == nil {
			t.Errorf("precondition failed: '%s' is not supposed to be a command in the test environment", nonExistentCommand)
			t.FailNow()
		}
		_, err := Run(Cmd{
			Command: nonExistentCommand,
		})
		if !errors.Is(err, ErrCommandNotFound) {
			t.Errorf("expected error '%v' to be an instance of %v\n", err, ErrCommandNotFound)
		}
	})

	t.Run("Empty command returns ErrCommandProcessFailed", func(t *testing.T) {
		_, err := Run(Cmd{})
		if !errors.Is(err, ErrCommandProcessFailed) {
			t.Errorf("expected error '%v' to be an instance of %v\n", err, ErrCommandProcessFailed)
		}
	})

	t.Run("Result contains value from stdout", func(t *testing.T) {
		expected := "passed"
		command := "/bin/sh"
		args := []string{"-c", fmt.Sprintf("echo %s", expected)}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "ECHO", expected}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		actual := strings.TrimSpace(res.Stdout.String())
		if actual != expected {
			t.Errorf("expected '%s'; got '%s'\n", expected, actual)
		}
	})

	t.Run("Result contains value from stderr", func(t *testing.T) {
		expected := "passed"
		command := "/bin/sh"
		args := []string{"-c", fmt.Sprintf("1>&2 echo %s", expected)}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "1>&2", "ECHO", expected}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		actual := strings.TrimSpace(res.Stderr.String())
		if actual != expected {
			t.Errorf("expected '%s'; got '%s'\n", expected, actual)
		}
	})

	t.Run("Result contains exit code", func(t *testing.T) {
		expected := 17
		command := "/bin/sh"
		args := []string{"-c", fmt.Sprintf("exit %d", expected)}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "exit", fmt.Sprintf("%d", expected)}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		if res.ExitCode != 17 {
			t.Errorf("expected %d; got %d\n", expected, res.ExitCode)
		}
	})

	t.Run("Command runs with specified Env", func(t *testing.T) {
		expected := "passed"
		env := []string{fmt.Sprintf("RESULT=%s", expected)}
		command := "/bin/sh"
		args := []string{"-c", "echo $RESULT"}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "ECHO", "%RESULT%"}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
			Env:     env,
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		actual := strings.TrimSpace(res.Stdout.String())
		if actual != expected {
			t.Errorf("expected '%s'; got '%s'\n", expected, actual)
		}
	})

	t.Run("Command reads stdin", func(t *testing.T) {
		expected, err := os.ReadFile("LICENSE")
		if err != nil {
			t.Errorf("precondition failed: failed to read LICENSE file: %v", err)
			t.FailNow()
		}
		command := "/bin/sh"
		args := []string{"-c", "cat"}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "findstr", "x*"}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
			Stdin:   bytes.NewReader(expected),
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		actual := res.Stdout.Bytes()
		if bytes.Compare(expected, actual) != 0 {
			t.Errorf("expected '%s'; got '%s'\n", expected, actual)
		}
	})

	t.Run("Command runs in directory specified by Dir", func(t *testing.T) {
		expected, err := os.ReadFile("testdata/dir.txt")
		if err != nil {
			t.Errorf("precondition failed: failed to read testdata/dir.txt: %v", err)
			t.FailNow()
		}
		command := "/bin/sh"
		args := []string{"-c", "cat dir.txt"}
		if runtime.GOOS == "windows" {
			command = "cmd.exe"
			args = []string{"/c", "type", "dir.txt"}
		}
		res, err := Run(Cmd{
			Command: command,
			Args:    args,
			Dir:     "testdata",
		})
		if err != nil {
			t.Errorf("unexpected error: %v\n", err)
			t.FailNow()
		}
		actual := res.Stdout.Bytes()
		if bytes.Compare(expected, actual) != 0 {
			t.Errorf("expected '%s'; got '%s'\n", expected, actual)
		}
	})
}
