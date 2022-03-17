package functional

import "fmt"

// Chain will create a function that composes the provided functions
// right to left, such that
//  f(x) == Chain(f)(x)
//  f(g(x)) == Chain(f, g)(x)
//  f(g(h(x))) == Chain(f, g, h)(x)
// and so on.
//
// Providing a nil function in the function slice will cause Chain
// to panic; however, providing a nil slice of functions is allowed.
// In the latter case, invoking the returned function will simply
// return the input.
func Chain[T any](fns ...func(T) T) func(T) T {
	f := func(t T) T { return t }
	for _, fn := range fns {
		if fn == nil {
			bork("nil chain function")
		}

		old := f
		fnCopy := fn
		f = func(t T) T { return old(fnCopy(t)) }
	}

	return f
}

// Compose will return a function that composes f and g; i.e.:
//  f(g(x)) == Compose(f, g)(x)
//
// If either function provided is nil, Compose will panic.
func Compose[In, Middle, Out any](f func(Middle) Out, g func(In) Middle) func(In) Out {
	switch {
	case f == nil:
		bork("nil f")
	case g == nil:
		bork("nil g")
	}

	return func(x In) Out {
		return f(g(x))
	}
}

// Filter will apply the provided filter to each element in
// the given list, returning the list of elements for which
// the given filter returned true.
//
// The returned list is guaranteed to be non-nil.
// Providing a nil slice is acceptable.
// The filter cannot be nil or Filter will panic.
func Filter[T any](list []T, filter func(T) bool) []T {
	if filter == nil {
		bork("nil filter")
	}

	filteredList := make([]T, 0, len(list))
	for _, v := range list {
		if filter(v) {
			filteredList = append(filteredList, v)
		}
	}

	return filteredList
}

// Map will map all values in the given list via the provided
// mapper. The returned list will contain the results of
// invoking the mapper on each element of the original list.
//
// The returned list is guaranteed to be non-nil.
// Providing a nil slice is acceptable.
// The mapper cannot be nil or Map will panic.
func Map[From, To any](list []From, mapper func(From) To) []To {
	if mapper == nil {
		bork("nil mapper")
	}

	mappedList := make([]To, 0, len(list))
	for _, v := range list {
		mappedList = append(mappedList, mapper(v))
	}

	return mappedList
}

// Reduce will apply the provided reducer to the given list,
// storing the results of each iteration's invocation of r,
// then returning the final value.
//
// Providing a nil slice is acceptable.
// The reducer cannot be nil or Reduce will panic.
func Reduce[From, To any](list []From, reducer func(accum To, cur From) To) To {
	if reducer == nil {
		bork("nil reducer")
	}

	var accum To
	for _, v := range list {
		accum = reducer(accum, v)
	}

	return accum
}

// bork will call panic with the prefix "functional".
func bork(v interface{}) {
	panic(fmt.Sprintf("functional: %+v", v))
}
