package exec

import (
	"fmt"
	"os/exec"
)

// CommandOutputFunc is a func which takes a command name and an optional
// number of args and produces a byte slice of output or an error.
type CommandOutputFunc func(cmd Cmd) ([]byte, error)

// CommandRunFunc is a func which takes a command name and an optional
// number of args and runs it, returning any errors.
type CommandRunFunc func(cmd Cmd) error

var (
	// commandOutputFn is pointing to the function that will be called by
	// CommandOutput. Can be overridden using FakeCommandOutput.
	commandOutputFn = commandOutput

	// commandRunFn is pointing to the function that will be called by
	// CommandRun. Can be overridden using FakeCommandRun.
	commandRunFn = commandRun
)

// CommandOutput runs the a command with given args and returns its standard
// output. Any returned error will usually be of type *ExitError.
//
// In the normal case, this just internally calls
// exec.Command(name, args...).Output() and returns the result.
//
// In tests the behaviour can be changed. See the documentation of the
// FakeCommandOutput func.
func CommandOutput(name string, args ...string) ([]byte, error) {
	return commandOutputFn(Cmd{name, args})
}

// commandOutput is a CommandOutputFunc which directly calls
// exec.Command(name, args...).Output() and returns the result.
func commandOutput(cmd Cmd) ([]byte, error) {
	output, err := exec.Command(cmd.Name, cmd.Args...).Output()

	return output, convertExitError(err)
}

// CommandRun runs the a command with given args. Any returned error will
// usually be of type *ExitError.
//
// In the normal case, this just internally calls
// exec.Command(name, args...).Run().
//
// In tests the behaviour can be changed. See the documentation of the
// FakeCommandRun func.
func CommandRun(name string, args ...string) error {
	return commandRunFn(Cmd{name, args})
}

// commandRun is a CommandRunFunc which directly calls
// exec.Command(name, args...).Run() and returns the result.
func commandRun(cmd Cmd) error {
	err := exec.Command(cmd.Name, cmd.Args...).Run()

	return convertExitError(err)
}

// FakeCommandOutput replaces all calls of CommandOutput with given fn in
// tests. The returned func must be called after the tests are finished to
// restore CommandOutput to avoid unexpected behaviour.
//
// 	 fakeCommandOutputFn := func(name string, args ...string) ([]byte, error) {
// 		 if name == "foo" && exec.ArgsMatch(args, "--bar", "baz") {
// 			 return []byte(`the output`), nil
// 		 }
//
// 		 return nil, &exec.ExitError{
// 			 ProcessState: &exec.FakeProcessState{ExitStatus: 1},
// 		 }
// 	 }
//
// 	 restore := FakeCommandOutput(fakeCommandOutputFn)
// 	 defer restore()
func FakeCommandOutput(fn CommandOutputFunc) func() {
	currentFn := commandOutputFn
	commandOutputFn = func(cmd Cmd) ([]byte, error) {
		output, err := fn(cmd)

		return output, maybeWrapExitError(err)
	}

	return func() { commandOutputFn = currentFn }
}

// FakeCommandRun replaces all calls of CommandRun with given fn in tests. The
// returned func must be called after the tests are finished to restore
// CommandRun to avoid unexpected behaviour.
//
// 	 fakeCommandRunFn := func(name string, args ...string) error {
// 		 if name == "foo" && exec.ArgsMatch(args, "--bar", "baz") {
// 			 return nil
// 		 }
//
// 		 return &exec.ExitError{
// 			 ProcessState: &exec.FakeProcessState{ExitStatus: 1},
// 		 }
// 	 }
//
// 	 restore := FakeCommandRun(fakeCommandRunFn)
// 	 defer restore()
func FakeCommandRun(fn CommandRunFunc) func() {
	currentFn := commandRunFn
	commandRunFn = func(cmd Cmd) error {
		return maybeWrapExitError(fn(cmd))
	}

	return func() { commandRunFn = currentFn }
}

// Cmd is a container type which wrap the command name and args. It has some
// methods attached to it which help matching commands in tests.
type Cmd struct {
	Name string
	Args []string
}

// ArgsMatch returns true if the command's args match the provided ones
// exactly.
func (c Cmd) ArgsMatch(args ...string) bool {
	if len(c.Args) != len(args) {
		return false
	}

	for i, arg := range c.Args {
		if arg != args[i] {
			return false
		}
	}

	return true
}

// Matches returns true if the command's name and args match the provided ones
// exactly.
func (c Cmd) Matches(name string, args ...string) bool {
	if name != c.Name {
		return false
	}

	return c.ArgsMatch(args...)
}

// Output internally calls exec.Command(c.Name, c.Args...).Output() and returns
// the result. This can be used to execute certain commands in tests while
// mocking others.
func (c Cmd) Output() ([]byte, error) {
	return commandOutput(c)
}

// Run internally calls exec.Command(c.Name, c.Args...).Run() and returns
// potential error. This can be used to execute certain commands in tests while
// mocking others.
func (c Cmd) Run() error {
	return commandRun(c)
}

// ProcessState is the interface satisfied by os.ProcessState containing only
// methods that might be usful in tests. This is used by ExitError so the
// ProcessState can be mocked in tests, e.g. if it is required to test
// behaviour based on different exit codes.
type ProcessState interface {
	String() string
	ExitCode() int
}

// ExitError has the same structure as exec.ExitError with the only difference
// that the embedded ProcessState is an interface instead of *os.ProcessState
// so that it can be faked in tests.
type ExitError struct {
	ProcessState
	Stderr []byte
}

// Error implements error.
func (e *ExitError) Error() string {
	return e.ProcessState.String()
}

// FakeProcessState is a fake process state which lets users fake the exit
// status of a command in tests.
type FakeProcessState struct {
	ExitStatus   int
	ErrorMessage string
}

// ExitCode implements ProcessState.
func (p *FakeProcessState) ExitCode() int {
	return p.ExitStatus
}

// String implements fmt.Stringer.
func (p *FakeProcessState) String() string {
	if len(p.ErrorMessage) > 0 {
		return p.ErrorMessage
	}

	return fmt.Sprintf("exit status %d", p.ExitStatus)
}

func maybeWrapExitError(err error) error {
	_, ok := err.(*ExitError)
	if ok || err == nil {
		return err
	}

	return &ExitError{
		ProcessState: &FakeProcessState{
			ExitStatus:   1,
			ErrorMessage: err.Error(),
		},
	}
}

func convertExitError(err error) error {
	exitError, ok := err.(*exec.ExitError)
	if err == nil || !ok {
		return err
	}

	return &ExitError{
		ProcessState: exitError.ProcessState,
		Stderr:       exitError.Stderr,
	}
}
