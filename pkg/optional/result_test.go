package optional_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/standoffvenus/functional/v2/pkg/optional"
	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	const Value = 42
	r := optional.Ok(Value)
	assert.True(t, r.Ok())
	assert.NoError(t, r.Err())
	assert.Equal(t, Value, r.Get())
	assert.Equal(t, Value, r.Expect())
}

func TestErr(t *testing.T) {
	var Error error = errors.New("error")
	r := optional.Err[int](Error)
	assert.False(t, r.Ok())
	assert.ErrorIs(t, r.Err(), Error)
	assert.Equal(t, int(0), r.Get())
	assert.Panics(t, func() { r.Expect() })
}

func TestResultStringWithValue(t *testing.T) {
	const Value = 42
	r := optional.Ok(Value)
	assert.Equal(t, strconv.FormatInt(Value, 10), r.String())
}

func TestResultStringWithError(t *testing.T) {
	var Error error = errors.New("message")
	r := optional.Err[int](Error)
	assert.Equal(t, Error.Error(), r.String())
}
