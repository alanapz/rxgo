package utils

import (
	"reflect"
	"slices"
	"sync"
)

type Never interface {
	unimplementable()
}

func Of[T any](values ...T) []T {
	return values
}

func Append[T any](slice *[]T, values ...T) {
	*slice = append(*slice, values...)
}

func Assert(condition bool) {
	if !condition {
		panic("Assertion failed")
	}
}

var nullableTypes = Of(reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice)

func Require(values ...any) {
	for value := range slices.Values(values) {
		reflectedValue := reflect.ValueOf(value)
		if slices.Contains(nullableTypes, reflectedValue.Kind()) && reflectedValue.IsNil() {
			panic("expected")
		}
	}
}

func Coalesce[T any](values ...T) T {

	for value := range slices.Values(values) {

		if !reflect.ValueOf(value).IsNil() {
			return value
		}
	}

	return Zero[T]()
}

func MustCoalesce[T any](values ...T) T {

	for value := range slices.Values(values) {

		if !reflect.ValueOf(value).IsNil() {
			return value
		}
	}

	panic("No values supplied")
}

func Ternary[T any](condition bool, ifTrue, ifFalse T) T {
	if condition {
		return ifTrue
	} else {
		return ifFalse

	}
}

func Zero[T any]() T {
	var value T
	return value
}

const mutexLocked = 1
const mutexUnlocked = 0

func AssertLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("mu").FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func AssertUnlocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("mu").FieldByName("state")
	return state.Int()&mutexUnlocked == mutexUnlocked
}
