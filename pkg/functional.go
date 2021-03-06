package functional

import (
	"sort"

	"github.com/standoffvenus/functional/v2/pkg/iterator"
)

// Break is a function that should be called when the caller
// wishes to break from a loop.
type Break = func()

// Comparable represents a type that can be compared to another
// instance of itself.
type Comparable interface {
	// Less should return true if the argument is considered
	// "greater than" the instance, i.e.
	//  var a, b Comparable
	//  a.Less(b) // == a < b
	Less(Comparable) bool
}

type (
	// comparables is used implement sort.Interface on a collection
	// of generic T.
	comparables[T Comparable] []T

	// sorted is used to represent a sorted iterator. It is returned
	// from Sort to quickly check if an iterator is sorted.
	sorted[T any] struct{ iterator.Iterator[T] }
)

// All will return whether the provided function holds true over
// all values in the iterator. If the iterator is empty, All will
// return true. All short-curcuits on the first value "x" such
// that fn(x) == false.
func All[T any](iter iterator.Iterator[T], fn func(T) bool) bool {
	return !Any(iter, func(t T) bool { return !fn(t) })
}

// Any will return whether the provided function holds true for
// any value in the iterator. If the iterator is empty, Any returns
// false. Any short-curcuits on the first value "x" such that
// fn(x) == true.
func Any[T any](iter iterator.Iterator[T], fn func(T) bool) bool {
	any := false
	ForEach(iter, func(t T, stop Break) {
		if any = fn(t); any {
			stop()
		}
	})

	return any
}

// Collect will call Next(), storing the results in a slice
// until None is encountered.
func Collect[T any](iter iterator.Iterator[T]) []T {
	slice := allocate[T](iter)
	ForEach(iter, func(t T, b Break) {
		slice = append(slice, t)
	})

	return slice
}

// CollectToChan will call Next(), sending the results to the
// returned channel on a separate Goroutine until None is
// encountered.
func CollectToChan[T any](iter iterator.Iterator[T]) <-chan T {
	ch := make(chan T, getSizeHint(iter))
	go func(c chan T) {
		defer close(c)
		ForEach(iter, func(t T, _ Break) {
			c <- t
		})
	}(ch)

	return ch
}

// Equal will check if two iterators equal by collecting their
// values and comparing the resulting slices. If the iterator's
// are different sizes, false is returned.
func Equal[T comparable](a, b iterator.Iterator[T]) bool {
	// Preliminary check on length to avoid collecting
	// both iterators if possible
	if getSizeHint(a) != getSizeHint(b) {
		return false
	}

	aValues, bValues := Collect(a), Collect(b)
	if len(aValues) != len(bValues) {
		return false
	}

	for idx := 0; idx < len(aValues); idx++ {
		if aValues[idx] != bValues[idx] {
			return false
		}
	}

	return true
}

// Filter will return an iterator with every value "x" in
// the given iterator such that fn(x) holds true.
func Filter[T any](iter iterator.Iterator[T], fn func(T) bool) iterator.Iterator[T] {
	filtered := iterator.Slice[T]{Values: allocate[T](iter)}
	ForEach(iter, func(t T, _ Break) {
		if fn(t) {
			filtered.Values = append(filtered.Values, t)
		}
	})

	return &filtered
}

// ForEach will call the provided function with each element
// returned from Next(), stopping iteration once None is returned.
// To break out of execution early, invoke Break.
func ForEach[T any](iter iterator.Iterator[T], fn func(T, Break)) {
	if iter == nil {
		return
	}

	loop := true
	stop := func() { loop = false }

	for loop {
		if opt := iter.Next(); opt.IsSome() {
			fn(opt.Expect(), stop)
		} else {
			stop()
		}
	}
}

// Map will return an iterator containing the results of
// invoking fn for each value of the provided iterator.
func Map[From, To any](iter iterator.Iterator[From], fn func(From) To) iterator.Iterator[To] {
	mapped := iterator.Slice[To]{Values: allocate[To](iter)}
	ForEach(iter, func(x From, _ Break) {
		mapped.Values = append(mapped.Values, fn(x))
	})

	return &mapped
}

// Reduce will invoke the provided function on each element
// of the given iterator, assigning a temporary variable to
// the results of each invocation, before returning the final
// value.
//
// The first argument passed to fn will be the current
// "accumulated" value from previous invocations, whereas the
// second argument will be the most recent result of calling
// iter.Next().
func Reduce[From, To any](iter iterator.Iterator[From], fn func(accum To, cur From) To) To {
	var accumulator To
	ForEach(iter, func(x From, _ Break) {
		accumulator = fn(accumulator, x)
	})

	return accumulator
}

// Sort will sort the provided iterator if it is not already sorted.
// If stable is set to true, the iterator will be sorted via sort.Stable.
// Otherwise, sort.Sort will be used.
//
// Sort's algorithm is as follows:
//  - If iter is the result of calling Sort, return iter (already sorted)
//  - Otherwise, collect all values from the provided iterator into a slice
//  - Check if the collected slice is sorted via sort.IsSorted
//  - If unsorted, sort the slice
func Sort[T Comparable](iter iterator.Iterator[T], stable bool) iterator.Iterator[T] {
	if _, ok := iter.(sorted[T]); ok {
		return iter
	}

	values := comparables[T](Collect(iter))
	if !sort.IsSorted(values) {
		if stable {
			sort.Stable(values)
		} else {
			sort.Sort(values)
		}
	}

	return sorted[T]{&iterator.Slice[T]{Values: []T(values)}}
}

// allocate will allocate a slice with some backing memory (not
// zeroed) equal to the size of the provided iterator's count
// if the iterator implements Enumerable.
func allocate[T, Source any](iter iterator.Iterator[Source]) []T {
	return make([]T, 0, getSizeHint(iter))
}

// getSizeHint will return iter.Count() if iter implements
// Enumerable. Otherwise, getSizedHint will return a default.
func getSizeHint[T any](iter iterator.Iterator[T]) int {
	const defaultSize = 16
	if sized, ok := iter.(iterator.Enumerable[T]); ok {
		if count := sized.Count(); count > 0 {
			return count
		}
	}

	return defaultSize
}

func (array comparables[T]) Len() int {
	return len(array)
}

func (array comparables[T]) Less(i, j int) bool {
	return array[i].Less(array[j])
}

func (array comparables[T]) Swap(i, j int) {
	tmp := array[i]
	array[i] = array[j]
	array[j] = tmp
}
