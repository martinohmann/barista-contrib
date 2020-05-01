package modules

import (
	"errors"
	"testing"

	"barista.run/bar"
	"barista.run/modules/static"
	"barista.run/outputs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := NewRegistry()

	r.Add(static.New(outputs.Text("")))

	require.NoError(t, r.Err())
	assert.Len(t, r.Modules(), 1)

	r.Add(nil)

	require.NoError(t, r.Err())
	assert.Len(t, r.Modules(), 1)

	r.Addf(func() (bar.Module, error) {
		return nil, nil
	})

	require.NoError(t, r.Err())
	assert.Len(t, r.Modules(), 1)

	r.Addf(func() (bar.Module, error) {
		return static.New(outputs.Text("")), nil
	})

	require.NoError(t, r.Err())
	assert.Len(t, r.Modules(), 2)

	moduleErr := errors.New("error")

	r.Addf(func() (bar.Module, error) {
		return nil, moduleErr
	})

	require.Error(t, r.Err())
	assert.Len(t, r.Modules(), 2)

	r.Add(static.New(outputs.Text("")))

	require.Error(t, r.Err())
	require.Same(t, moduleErr, r.Err())
	assert.Len(t, r.Modules(), 2)

	r.Addf(func() (bar.Module, error) {
		return static.New(outputs.Text("")), nil
	})

	require.Error(t, r.Err())
	require.Same(t, moduleErr, r.Err())
	assert.Len(t, r.Modules(), 2)
}
