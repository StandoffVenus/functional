package functional_test

import (
	"testing"

	functional "github.com/standoffvenus/functional/v2/pkg"
	"github.com/standoffvenus/functional/v2/pkg/iterator"
	"github.com/standoffvenus/functional/v2/pkg/optional"
	"github.com/stretchr/testify/assert"
)

func GreaterThan0(x int) bool { return x > 0 }

func TestAllWithAllTrue(t *testing.T) {
	iter := Iterator(1, 2, 3)
	assert.True(t, functional.All(iter, GreaterThan0))
}

func TestAllWithFalseValue(t *testing.T) {
	iter := Iterator(3, 5, 0)
	assert.False(t, functional.All(iter, GreaterThan0))
}

func TestAllWithNoValues(t *testing.T) {
	assert.True(t, functional.All(Iterator[int](), GreaterThan0))
}

func TestAnyWithAllFalse(t *testing.T) {
	iter := Iterator(-2, -1, 0)
	assert.False(t, functional.Any(iter, GreaterThan0))
}

func TestAnyWithTrueValue(t *testing.T) {
	iter := Iterator(-2, 0, 2)
	assert.True(t, functional.Any(iter, GreaterThan0))
}

func TestAnyWithNoValues(t *testing.T) {
	assert.False(t, functional.Any(Iterator[int](), GreaterThan0))
}

func TestCollect(t *testing.T) {
	ints := []int{1, 2, 3}
	iter := &iterator.Slice[int]{Values: ints}
	collected := functional.Collect[int](iter)

	assert.Equal(t, ints, collected)
}

func TestCollectToChan(t *testing.T) {
	ints := []int{1, 2, 3}
	iter := &iterator.Slice[int]{Values: ints}
	collected := functional.CollectToChan[int](iter)

	AssertEqualChan(t, ints, collected)
}

func TestCollectToChanNoDeadlock(t *testing.T) {
	const Value = 42
	f := func() optional.Option[int] { return optional.Some(Value) }
	iter := iterator.Func[int](f)                    // Func iterator doesn't have size hint,
	collected := functional.CollectToChan[int](iter) // so this channel is unbuffered.

	assert.Equal(t, Value, <-collected)
}

func TestFilter(t *testing.T) {
	ints := []int{-1, 0, 1}
	iter := &iterator.Slice[int]{Values: ints}
	filtered := functional.Filter[int](iter, GreaterThan0)

	AssertIteratorEqual(t, []int{1}, filtered)
}

func TestForEach(t *testing.T) {
	ints := []int{-1, 0, 1}
	iter := &iterator.Slice[int]{Values: ints}
	loopedValues := make([]int, 0, iter.Count())

	functional.ForEach[int](iter, func(x int, _ functional.Break) {
		loopedValues = append(loopedValues, x)
	})

	assert.Equal(t, ints, loopedValues)
}

func TestForEachNilIterator(t *testing.T) {
	assert.NotPanics(t, func() {
		functional.ForEach(nil, func(_ int, _ functional.Break) {})
	})
}

func TestForEachCanBreak(t *testing.T) {
	ints := []int{-1, 0, 1}
	iter := &iterator.Slice[int]{Values: ints}
	loopedValues := make([]int, 0, iter.Count())

	functional.ForEach[int](iter, func(x int, stop functional.Break) {
		loopedValues = append(loopedValues, x)
		stop()
	})

	assert.Less(t, len(loopedValues), len(ints))
	assert.Subset(t, ints, loopedValues)
}

func TestMap(t *testing.T) {
	ints := []int{0, 1, 2}
	iter := &iterator.Slice[int]{Values: ints}
	expected := []int{0, 1, 4}

	mapped := functional.Map[int](iter, func(x int) int { return x * x })

	AssertIteratorEqual(t, expected, mapped)
}

func TestMapToDifferentType(t *testing.T) {
	ints := []int{0, 1, 2}
	iter := &iterator.Slice[int]{Values: ints}
	expected := []float32{0, 1, 2}

	mapped := functional.Map[int](iter, func(x int) float32 { return float32(x) })

	AssertIteratorEqual(t, expected, mapped)
}

func TestReduce(t *testing.T) {
	ints := []int{0, 1, 2}
	iter := &iterator.Slice[int]{Values: ints}
	expected := 0
	for _, i := range ints {
		expected += i * i
	}

	reduced := functional.Reduce[int](iter, func(accum int, cur int) int { return accum + cur*cur })

	assert.Equal(t, expected, reduced)
}

func TestReduceToDifferentType(t *testing.T) {
	ints := []int{0, 1, 2}
	iter := &iterator.Slice[int]{Values: ints}
	expected := 0.0
	for _, i := range ints {
		expected += float64(i)
	}

	reduced := functional.Reduce[int](iter, func(accum float64, cur int) float64 { return accum + float64(cur) })

	assert.Equal(t, expected, reduced)
}

func AssertIteratorEqual[T comparable](t *testing.T, expected []T, iter iterator.Iterator[T]) bool {
	for idx, v := range expected {
		if v != iter.Next().Expect() {
			t.Errorf("expected[%d] was not equal to next iterator value (%v)", idx, v)
			return false
		}
	}

	return true
}

func AssertEqualChan[T any](t *testing.T, expected []T, ch <-chan T) bool {
	slice := make([]T, 0, len(ch))
	for v := range ch {
		slice = append(slice, v)
	}

	return assert.Equal(t, expected, slice)
}

func Iterator[T any](values ...T) iterator.Iterator[T] {
	return &iterator.Slice[T]{Values: values}
}
