package optional

// Ok will return an OK result with the given
// value.
func Ok[T any](t T) Result[T] {
	return Result[T]{opt: Some(t)}
}

// Err will return an error result with the
// given error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Result represents an optional value whose
// absence represents an error.
//
// A zero-value result is considered erroneous
// for the sake of convenience - it will not
// contain an error. Prefer to copy results.
type Result[T any] struct {
	opt Option[T]
	err error
}

// Ok will return whether the result is considered
// erroneous.
func (r Result[T]) Ok() bool {
	return r.opt.IsSome()
}

// Get will return the value stored in the result.
func (r Result[T]) Get() T {
	return r.opt.Get()
}

// Err will return the error stored in the result.
func (r Result[T]) Err() error {
	return r.err
}

// Expect is the same as Get but panics if the result
// is not considered OK; in other words, if r.Ok returns
// false, Expect panics.
func (r Result[T]) Expect() T {
	if !r.Ok() {
		panic("optional: Expect() called on error result")
	}

	return r.opt.value
}

// String will return the result's value formatted using fmt.Sprintf,
// or the error string if the result is erroneous.
func (r Result[T]) String() string {
	if r.Ok() {
		return r.opt.String()
	}

	return r.err.Error()
}
