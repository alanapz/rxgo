package rx

import (
	"reflect"
	"slices"
	"sync"
)

var NoOp = func() {

}

type Never interface {
	unimplementable()
}

func Append[T any](slice *[]T, values ...T) {
	*slice = append(*slice, values...)
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

func AssertLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("mu").FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func Cast[T any](value reflect.Value, present bool) T {
	if !present {
		return Zero[T]()
	}
	return value.Interface().(T)
}
