package optional_test

import (
	"strconv"
	"testing"

	"github.com/standoffvenus/functional/v2/pkg/optional"
	"github.com/stretchr/testify/assert"
)

func TestSome(t *testing.T) {
	const Value = 42
	v := optional.Some(Value)
	assert.True(t, v.IsSome())
	assert.Equal(t, Value, v.Get())
	assert.Equal(t, Value, v.Expect())
}

func TestNone(t *testing.T) {
	v := optional.None[int]()
	assert.False(t, v.IsSome())
	assert.Equal(t, int(0), v.Get())
	assert.Panics(t, func() { v.Expect() })
}

func TestOptionZeroIsNone(t *testing.T) {
	assert.False(t, optional.Option[int]{}.IsSome())
}

func TestOptionStringWithNoValue(t *testing.T) {
	v := optional.None[int]()
	assert.Equal(t, "None", v.String())
}

func TestOptionStringWithValue(t *testing.T) {
	const Value = 42
	v := optional.Some(Value)
	assert.Equal(t, strconv.FormatInt(Value, 10), v.String())
}
