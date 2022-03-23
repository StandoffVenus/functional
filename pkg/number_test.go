package functional_test

import (
	"math"
	"testing"
	"testing/quick"

	functional "github.com/standoffvenus/functional/v2/pkg"
	"github.com/standoffvenus/functional/v2/pkg/iterator"
	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	quick.Check(
		func(floats []float64) bool {
			iter := &iterator.Slice[float64]{Values: floats}
			expectedSum := float64(0)
			for _, v := range iter.Values {
				expectedSum += v
			}

			return expectedSum == functional.Sum[float64](iter)
		},
		nil,
	)
}

func TestMultiplyScalar(t *testing.T) {
	const factor float64 = 2.5
	quick.Check(
		func(floats []float64) bool {
			iter := &iterator.Slice[float64]{Values: floats}
			expectedProduct := float64(0)
			for _, v := range iter.Values {
				expectedProduct *= v
			}

			return expectedProduct == functional.MultiplyScalar[float64](iter)
		},
		nil,
	)
}

func TestMultiplyVector(t *testing.T) {
	const factor float64 = 2.5
	quick.Check(
		func(floats []float64) bool {
			iter := &iterator.Slice[float64]{Values: floats}
			var expected []float64
			for _, v := range iter.Values {
				expected = append(expected, v*factor)
			}

			actualIterator := functional.MultiplyVector[float64](iter, factor)

			return AssertIteratorEqual[float64](t, expected, actualIterator)
		},
		nil,
	)
}

func TestDotProduct(t *testing.T) {
	a := &iterator.Slice[float64]{Values: []float64{6, -2, -1}}
	b := &iterator.Slice[float64]{Values: []float64{2, 10, 2}}

	assert.Equal(t, float64(-10), functional.DotProduct[float64](a, b))
}

func TestDotProductPanicsOnDifferentDimensions(t *testing.T) {
	assert.Panics(t, func() {
		a := &iterator.Slice[int]{}
		b := &iterator.Slice[int]{Values: []int{42}}

		functional.DotProduct[int](a, b)
	})
}

func TestSquare(t *testing.T) {
	iter := &iterator.Slice[float64]{Values: []float64{1, 2, 3, 4}}
	squaredIterator := functional.Square[float64](iter)

	var expected []float64
	for _, v := range iter.Values {
		expected = append(expected, v*v)
	}

	AssertIteratorEqual(t, expected, squaredIterator)
}

func TestTriple(t *testing.T) {
	iter := &iterator.Slice[float64]{Values: []float64{1, 2, 3, 4}}
	tripledIterator := functional.Triple[float64](iter)

	var expected []float64
	for _, v := range iter.Values {
		expected = append(expected, v*v*v)
	}

	AssertIteratorEqual(t, expected, tripledIterator)
}

func TestQuadruple(t *testing.T) {
	iter := &iterator.Slice[float64]{Values: []float64{1, 2, 3, 4}}
	quadrupledIterator := functional.Quadruple[float64](iter)

	var expected []float64
	for _, v := range iter.Values {
		expected = append(expected, v*v*v*v)
	}

	AssertIteratorEqual(t, expected, quadrupledIterator)
}

func TestToPower(t *testing.T) {
	const Power float64 = 9.6
	iter := &iterator.Slice[float64]{Values: []float64{1, 2, 3, 4}}
	toPowerIterator := functional.ToPower[float64](iter, Power)

	var expected []float64
	for _, v := range iter.Values {
		expected = append(expected, math.Pow(v, Power))
	}

	AssertIteratorEqual(t, expected, toPowerIterator)
}
