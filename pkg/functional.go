package functional

import "github.com/standoffvenus/functional/v2/pkg/iterator"

// Break is a function that should be called when the caller
// wishes to break from a loop.
type Break = func()

func All[T any](iter iterator.Iterator[T], fn func(T) bool) bool {
	return !Any(iter, func(t T) bool { return !fn(t) })
}

func Any[T any](iter iterator.Iterator[T], fn func(T) bool) bool {
	any := false
	ForEach(iter, func(t T, stop Break) {
		if any = fn(t); any {
			stop()
		}
	})

	return any
}

func Collect[T any](iter iterator.Iterator[T]) []T {
	slice := allocate[T](iter)
	ForEach(iter, func(t T, b Break) {
		slice = append(slice, t)
	})

	return slice
}

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

func Filter[T any](iter iterator.Iterator[T], fn func(T) bool) iterator.Iterator[T] {
	filtered := iterator.Slice[T]{Values: allocate[T](iter)}
	ForEach(iter, func(t T, _ Break) {
		if fn(t) {
			filtered.Values = append(filtered.Values, t)
		}
	})

	return &filtered
}

func ForEach[T any](iter iterator.Iterator[T], fn func(T, Break)) {
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

func Map[From, To any](iter iterator.Iterator[From], fn func(From) To) iterator.Iterator[To] {
	mapped := iterator.Slice[To]{Values: allocate[To](iter)}
	ForEach(iter, func(x From, _ Break) {
		mapped.Values = append(mapped.Values, fn(x))
	})

	return &mapped
}

func Reduce[From, To any](iter iterator.Iterator[From], fn func(accum To, cur From) To) To {
	var accumulator To
	ForEach(iter, func(x From, _ Break) {
		accumulator = fn(accumulator, x)
	})

	return accumulator
}

func allocate[T, Source any](iter iterator.Iterator[Source]) []T {
	return make([]T, 0, getSizeHint(iter))
}

func getSizeHint[T any](iter iterator.Iterator[T]) int {
	const defaultSize = 16
	if sized, ok := iter.(iterator.Enumerable[T]); ok {
		if count := sized.Count(); count > 0 {
			return count
		}
	}

	return defaultSize
}
