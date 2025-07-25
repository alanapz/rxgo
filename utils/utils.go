package utils

import (
	"reflect"
	"slices"
	"sync"
)

type Void = struct{}

type Never interface {
	unimplementable()
}

func Of[T any](values ...T) []T {
	return values
}

func Val[T any](value T) func() T {
	return func() T {
		return value
	}
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

func Require[T any](value T) T {
	reflectedValue := reflect.ValueOf(value)
	if slices.Contains(nullableTypes, reflectedValue.Kind()) && reflectedValue.IsNil() {
		panic("parameter required")
	}
	return value
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

func Panic(err error) {
	panic(err)
}

func DoNothing() {

}

const mutexLocked = 1
const mutexUnlocked = 0

func Lock(m *sync.Mutex) func() {
	m.Lock()
	return func() {
		defer m.Unlock()
	}
}

func AssertLocked(m *sync.Mutex) {
	state := reflect.ValueOf(m).Elem().FieldByName("mu").FieldByName("state")
	if state.Int()&mutexLocked != mutexLocked {
		panic("mutex was not held")
	}
}

func AssertUnlocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("mu").FieldByName("state")
	return state.Int()&mutexUnlocked == mutexUnlocked
}
