package errx

import (
	"errors"

	"golang.org/x/exp/constraints"
)

func AsError(err error) (xerr, bool) {
	var e xerr
	ok := errors.As(err, &e)
	return e, ok
}

func maxValue[T constraints.Ordered](collection []T) T {
	var value T

	if len(collection) == 0 {
		return value
	}

	value = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item > value {
			value = item
		}
	}

	return value
}

func minValue[T constraints.Ordered](collection []T) T {
	var value T

	if len(collection) == 0 {
		return value
	}

	value = collection[0]

	for i := 1; i < len(collection); i++ {
		item := collection[i]

		if item < value {
			value = item
		}
	}

	return value
}

// repeatBy builds a slice with values returned by N calls of callback.
func repeatBy[T any](count int, predicate func(index int) T) []T {
	result := make([]T, 0, count)

	for i := 0; i < count; i++ {
		result = append(result, predicate(i))
	}

	return result
}

// ternaryF is a 1 line if/else statement whose options are functions
// Play: https://go.dev/play/p/AO4VW20JoqM
func ternaryF[T any](condition bool, ifFunc func() T, elseFunc func() T) T {
	if condition {
		return ifFunc()
	}

	return elseFunc()
}

// coalesce returns the first non-empty arguments. Arguments must be comparable.
func coalesce[T comparable](v ...T) (result T, ok bool) {
	for _, e := range v {
		if e != result {
			result = e
			ok = true
			return
		}
	}

	return
}

func coalesceOrEmpty[T comparable](v ...T) T {
	result, _ := coalesce(v...)
	return result
}

// manipulateMap manipulates a slice and transforms it to a slice of another type.
func manipulateMap[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}
