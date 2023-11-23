package errx

import (
	"errors"
	"fmt"
)

func New(message string) error {
	return Errorf(message)
}

// Errorf formats an error and returns `errx.Error` object that satisfies `error`.
func Errorf(format string, args ...any) error {
	e := newError()
	e.msg = fmt.Sprintf(format, args...)
	return e
}

// Wrap wraps an error into an `errx.Error` object that satisfies `error`
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	return Wrapf(err, msg)
}

// Wrapf wraps an error into an `errx.Error` object that satisfies `error` and formats an error message.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	e := newError()
	e.err = err
	e.msg = fmt.Sprintf(format, args...)
	return e
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
