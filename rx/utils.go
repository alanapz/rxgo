package rx

import (
	"reflect"
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
