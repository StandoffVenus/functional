package optional_test

import (
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

func TestDefaultOptionIsNone(t *testing.T) {
	assert.False(t, optional.Option[int]{}.IsSome())
}
