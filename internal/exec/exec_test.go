package exec

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFakeCommandOutput(t *testing.T) {
	restore := FakeCommandOutput(func(cmd Cmd) ([]byte, error) {
		switch {
		case cmd.Matches("foo", "--bar", "baz"):
			return []byte(`foo`), nil
		case cmd.Matches("foo", "--error"):
			return nil, errors.New("some error")
		case cmd.Name == "bar":
			return []byte(`bar`), nil
		default:
			return nil, &ExitError{
				ProcessState: &FakeProcessState{
					ExitStatus: 42,
				},
			}
		}
	})
	defer restore()

	output, err := CommandOutput("foo", "--bar", "baz")
	require.NoError(t, err)
	assert.Equal(t, "foo", string(output))

	output, err = CommandOutput("foo", "--error")
	require.Error(t, err)

	exitError, ok := err.(*ExitError)
	require.True(t, ok)
	assert.Equal(t, 1, exitError.ExitCode())
	assert.Equal(t, "some error", exitError.Error())

	output, err = CommandOutput("bar")
	require.NoError(t, err)
	assert.Equal(t, "bar", string(output))

	output, err = CommandOutput("foo", "--somearg")
	require.Error(t, err)

	exitError, ok = err.(*ExitError)
	require.True(t, ok)
	assert.Equal(t, 42, exitError.ExitCode())
	assert.Equal(t, "exit status 42", exitError.Error())
}

func TestFakeCommandRun(t *testing.T) {
	restore := FakeCommandRun(func(cmd Cmd) error {
		switch {
		case cmd.Matches("foo", "--bar", "baz"):
			return nil
		case cmd.Matches("foo", "--error"):
			return errors.New("some error")
		default:
			return &ExitError{
				ProcessState: &FakeProcessState{
					ExitStatus: 42,
				},
			}
		}
	})
	defer restore()

	err := CommandRun("foo", "--bar", "baz")
	require.NoError(t, err)

	err = CommandRun("foo", "--error")
	require.Error(t, err)

	exitError, ok := err.(*ExitError)
	require.True(t, ok)
	assert.Equal(t, 1, exitError.ExitCode())
	assert.Equal(t, "some error", exitError.Error())

	err = CommandRun("someothercommand")
	require.Error(t, err)

	exitError, ok = err.(*ExitError)
	require.True(t, ok)
	assert.Equal(t, 42, exitError.ExitCode())
	assert.Equal(t, "exit status 42", exitError.Error())
}

func TestConvertExitError(t *testing.T) {
	assert.Nil(t, convertExitError(nil))
	assert.Equal(t, errors.New("foo"), convertExitError(errors.New("foo")))

	exitError := &exec.ExitError{
		ProcessState: &os.ProcessState{},
		Stderr:       []byte(`error`),
	}

	expectedError := &ExitError{
		ProcessState: exitError.ProcessState,
		Stderr:       []byte(`error`),
	}

	assert.Equal(t, expectedError, convertExitError(exitError))
}
