package iterator

import (
	"context"

	"github.com/standoffvenus/functional/v2/pkg/optional"
)

// Iterator represents a basic iterator on type T.
type Iterator[T any] interface {
	// Next will retrieve the next value in the iterator.
	// Typically, if None is returned, the iterator can
	// be considered exhausted.
	Next() optional.Option[T]
}

// BlockingIterator represents an iterator that may
// block indefinitely on its Next().
type BlockingIterator[T any] interface {
	// WaitForNext will block until another value
	// can be retrieved from the iterator. If the
	// context is canceled, WaitForNext will stop
	// blocking and None will be returned.
	WaitForNext(ctx context.Context) optional.Option[T]
}

// Enumerable represents an iterator with a size.
type Enumerable[T any] interface {
	Iterator[T]

	// Count returns the size of the iterator.
	Count() int
}

// Slice represents an iterator on a generic slice.
type Slice[T any] struct {
	// Values holds the slice to iterate. The slice
	// should not be modified after creating the iterator.
	//
	// A nil slice is equivalent to an exhausted iterator.
	Values []T

	index int
}

// Chan represents an iterator on a generic channel.
// A channel iterator will block forever when Next
// is called.
type Chan[T any] <-chan T

// Func represents an iterator on a generic function
// returning optional values. A nil function iterator
// is equivalent to an exhausted iterator.
type Func[T any] func() optional.Option[T]

var _ Iterator[int] = new(Slice[int])
var _ Iterator[int] = Chan[int](nil)
var _ Iterator[int] = Func[int](nil)

var _ BlockingIterator[int] = new(Slice[int])
var _ BlockingIterator[int] = Chan[int](nil)

var _ Enumerable[int] = new(Slice[int])

// Send will create a buffered channel, send all the provided
// values on it, then return the channel to the caller. Useful
// when a channel iterator is needed from a collection of values.
func Send[T any](values ...T) chan T {
	ch := make(chan T, len(values))
	for _, v := range values {
		ch <- v
	}

	return ch
}

// WaitForNext is a general-purpose method that simplifies waiting
// on Next() to return a value. If the provided context is canceled,
// WaitForNext returns None.
//
// If the provided iterator implements BlockingIterator, the type's
// WaitForNext method will be called directly. Otherwise, a Goroutine
// is started that will send the result of Next() to a buffered channel.
// As such, if the iterator blocks indefinitely, the Goroutine will be
// "leaked".
func WaitForNext[T any](ctx context.Context, iter Iterator[T]) optional.Option[T] {
	if blockingIter, ok := iter.(BlockingIterator[T]); ok {
		return blockingIter.WaitForNext(ctx)
	}

	return waitForNext(ctx, iter)
}

// WaitForNext will return the next value of the iterator
// or None.
//
// WaitForNext will never block since Next() is non-blocking.
// WaitForNext is implemented on Slice only to avoid
// allocating a Goroutine when passed to iterator.WaitForNext.
func (s *Slice[T]) WaitForNext(_ context.Context) optional.Option[T] { return s.Next() }

// Count will return the remaining number of elements to
// iterate.
func (s *Slice[T]) Count() int { return len(s.Values) - s.index }

// Next will return the first value of the underlying slice
// if there is one, advancing the
func (s *Slice[T]) Next() optional.Option[T] {
	if len(s.Values) > s.index {
		s.index++
		return optional.Some(s.Values[s.index-1])
	}

	return optional.None[T]()
}

// Next returns the result of waiting for the next value from the channel.
// If the channel is closed, None is returned.
//
// If the channel is nil, Next will block forever.
func (c Chan[T]) Next() optional.Option[T] {
	if v, ok := <-c; ok {
		return optional.Some(v)
	}

	return optional.None[T]()
}

// WaitForNext will wait until either the channel receives a value, it
// closes, or the provided context is canceled. For the latter two
// cases, None is returned.
func (c Chan[T]) WaitForNext(ctx context.Context) optional.Option[T] {
	select {
	case v, ok := <-c:
		if ok {
			return optional.Some(v)
		}
	case <-ctx.Done():
	}

	return optional.None[T]()
}

// Next will return the result of calling f if it is not nil.
// Otherwise, None is always returned.
func (f Func[T]) Next() optional.Option[T] {
	if f != nil {
		return f()
	}

	return optional.None[T]()
}

func waitForNext[T any](ctx context.Context, iter Iterator[T]) optional.Option[T] {
	ch := make(chan optional.Option[T], 1)
	go func() {
		defer close(ch)
		ch <- iter.Next()
	}()

	select {
	case o := <-ch:
		return o
	case <-ctx.Done():
	}

	return optional.None[T]()
}
