package functional

import (
	"math"

	"github.com/standoffvenus/functional/v2/pkg/iterator"
)

// Number represents all numeric types in Go.
type Number interface {
	Rational | ~complex64 | ~complex128
}

// Rational represents all non-complex types in Go.
type Rational interface {
	~int8 | ~int16 | ~int32 | ~int | ~int64 |
		~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64 |
		~float32 | ~float64
}

// Sum will sum the elements of a numeric iterator.
func Sum[T Number](iter iterator.Iterator[T]) T {
	return Reduce(iter, func(accum, cur T) T { return accum + cur })
}

// MultiplyScalar will multiply all the elements of a
// numeric iterator together to produce their product.
func MultiplyScalar[T Number](iter iterator.Iterator[T]) T {
	return Reduce(iter, func(accum, cur T) T { return accum * cur })
}

// MultiplyVector will multiply all elements in the iterator
// by the provided factor, returning an iterator containing
// the products.
func MultiplyVector[T Number](iter iterator.Iterator[T], factor T) iterator.Iterator[T] {
	return Map(iter, func(x T) T { return x * factor })
}

// DotProduct will multiply each value of both iterators
// and return the sum of their products. If the iterators
// are different sizes, DotProduct will panic.
func DotProduct[T Number](a, b iterator.Enumerable[T]) T {
	if a.Count() != b.Count() {
		panic("functional: dot product on iterators with different dimensions")
	}

	return Reduce[T](a, func(accum T, x T) T {
		return accum + (x * b.Next().Expect())
	})
}

// Square will square each value in the iterator, returning
// an iterator containing the squares.
func Square[T Number](iter iterator.Iterator[T]) iterator.Iterator[T] {
	return Map(iter, func(x T) T { return x * x })
}

// Triple will triple each value in the iterator, returning
// an iterator containing the triples.
func Triple[T Number](iter iterator.Iterator[T]) iterator.Iterator[T] {
	return Map(iter, func(x T) T { return x * x * x })
}

// Quadruple will quadruple each value in the iterator,
// returning an iterator containing the quadruples.
func Quadruple[T Number](iter iterator.Iterator[T]) iterator.Iterator[T] {
	return Map(iter, func(x T) T { return x * x * x * x })
}

// ToPower will raise each element in the iterator to the
// provided power, returning the results in an iterator.
// Prefer Square, Triple, and Quadruple - ToPower uses
// math.Pow, which may be significantly less performant.
func ToPower[T Rational](iter iterator.Iterator[T], exp T) iterator.Iterator[T] {
	return Map(iter, func(x T) T { return T(math.Pow(float64(x), float64(exp))) })
}
