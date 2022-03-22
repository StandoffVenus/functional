// Package optional provides primitives for representing
// optional values.
package optional

import "fmt"

// Some constructs a new option with the provided
// value, representing the Some invariant.
func Some[T any](t T) Option[T] {
	return Option[T]{value: t, some: true}
}

// None constructs an empty option representing
// the None invariant.
func None[T any]() Option[T] {
	return Option[T]{}
}

// Option represents an optional value. If an
// Option does not have a value, it is referred
// to as "None". Likewise, an option with a
// value is "Some".
//
// A zero-value option is None and ready for use.
// Prefer to copy options.
type Option[T any] struct {
	some  bool
	value T
}

// IsSome returns true iff the option has
// a value.
func (o Option[T]) IsSome() bool {
	return o.some
}

// Get will retrieve the option's value.
// If None, the returned value is the zero
// value of T.
func (o Option[T]) Get() T {
	return o.value
}

// Expect is the same as Get, expect it
// panics if the option is None.
func (o Option[T]) Expect() T {
	if !o.IsSome() {
		panic("optional: Expect() called on None")
	}

	return o.value
}

// String will return the option's value
// formatted using fmt.Sprintf, or "None"
// if the option has no value.
func (o Option[T]) String() string {
	if o.IsSome() {
		return fmt.Sprintf("%+v", o.Expect())
	}

	return fmt.Sprintf("None")
}
