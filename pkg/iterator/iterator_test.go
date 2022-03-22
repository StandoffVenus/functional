package iterator_test

import (
	"context"
	"testing"

	"github.com/standoffvenus/functional/v2/pkg/iterator"
	"github.com/standoffvenus/functional/v2/pkg/optional"
	"github.com/stretchr/testify/assert"
)

var Values []int = []int{4, 9, 13}

func TestSliceCount(t *testing.T) {
	s := iterator.Slice[int]{
		Values: Values,
	}
	_ = s.Next()

	assert.Equal(t, len(s.Values)-1, s.Count())
}

func TestSliceNext(t *testing.T) {
	iter := &iterator.Slice[int]{
		Values: Values,
	}

	AssertIteratorMatches[int](t, iter, Values)
	AssertNextIsNone[int](t, iter)
}

func TestSliceWaitForNext(t *testing.T) {
	ctx := context.Background()
	iter := &iterator.Slice[int]{
		Values: Values,
	}

	AssertIteratorWaitForNextMatches[int](t, ctx, iter, Values)
	AssertWaitForNextIsNone[int](t, ctx, iter)
}

func TestChanNext(t *testing.T) {
	ch := iterator.Send(Values...)
	close(ch)
	iter := iterator.Chan[int](ch)

	AssertIteratorMatches[int](t, iter, Values)
	AssertNextIsNone[int](t, iter)
}

func TestChanWaitForNext(t *testing.T) {
	ctx := context.Background()
	ch := iterator.Send(Values...)
	close(ch)
	iter := iterator.Chan[int](ch)

	AssertIteratorWaitForNextMatches[int](t, ctx, iter, Values)
	AssertWaitForNextIsNone[int](t, ctx, iter)
}

func TestChanWaitForNextOnClosedChannel(t *testing.T) {
	ctx := context.Background()
	ch := make(chan int)
	close(ch)
	iter := iterator.Chan[int](ch)

	AssertWaitForNextIsNone[int](t, ctx, iter)
}

func TestChanWaitForNextOnCanceledContext(t *testing.T) {
	ctx := canceled()
	ch := make(chan int)
	defer close(ch)
	iter := iterator.Chan[int](ch)

	AssertWaitForNextIsNone[int](t, ctx, iter)
}

func TestFuncNext(t *testing.T) {
	iter := funcIteratorOf(Values)

	AssertIteratorMatches[int](t, iter, Values)
	AssertNextIsNone[int](t, iter)
}

func TestFuncNextOnNil(t *testing.T) {
	iter := iterator.Func[int](nil)

	AssertNextIsNone[int](t, iter)
}

func TestWaitForNext(t *testing.T) {
	ctx := context.Background()
	iter := funcIteratorOf(Values)
	expectedOption := optional.Some(Values[0])

	assert.Equal(t, expectedOption, iterator.WaitForNext[int](ctx, iter))
}

func TestWaitForNextOnBlockingIterator(t *testing.T) {
	ctx := canceled()
	iter := iterator.Chan[int](nil)

	assert.Equal(t, optional.None[int](), iterator.WaitForNext[int](ctx, iter))
}

func TestWaitForNextForeverBlocking(t *testing.T) {
	ctx := canceled()
	iter := funcIteratorOf[int](nil)

	assert.Equal(t, optional.None[int](), iterator.WaitForNext[int](ctx, iter))
}

func funcIteratorOf[T any](v []T) iterator.Func[T] {
	ch := iterator.Send(v...)
	defer close(ch)
	fn := func() optional.Option[T] { return iterator.Chan[T](ch).Next() }

	return iterator.Func[T](fn)
}

func AssertNextIsNone[T any](t *testing.T, iter iterator.Iterator[T]) bool {
	return assert.Equal(t, optional.None[T](), iter.Next())
}

func AssertWaitForNextIsNone[T any](t *testing.T, ctx context.Context, iter iterator.BlockingIterator[T]) bool {
	return assert.Equal(t, optional.None[T](), iter.WaitForNext(ctx))
}

func AssertIteratorMatches[T any](t *testing.T, iter iterator.Iterator[T], values []T) bool {
	for _, v := range values {
		if !assert.Equal(t, v, iter.Next().Expect()) {
			return false
		}
	}

	return true
}

func AssertIteratorWaitForNextMatches[T any](
	t *testing.T,
	ctx context.Context,
	iter iterator.BlockingIterator[T],
	values []T,
) bool {
	for _, v := range values {
		if !assert.Equal(t, v, iter.WaitForNext(ctx).Expect()) {
			return false
		}
	}

	return true
}

func canceled() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	return ctx
}
