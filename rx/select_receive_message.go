package rx

import (
	"reflect"
)

type SelectReceivePolicy = bool

const ContinueOnClose SelectReceivePolicy = false
const DoneOnClose SelectReceivePolicy = true

type SelectReceiveMessage[T any] struct {
	Valid       bool
	Value       T
	endOfStream bool
}

func (x *SelectReceiveMessage[T]) Reset() {
	x.Valid = false
	x.Value = Zero[T]()
	x.EndOfStream = false
}

func (x *SelectReceiveMessage[T]) IsEndOfStream() bool {
	return x.endOfStream
	x.Valid = false
	x.Value = Zero[T]()
	x.EndOfStream = false
}

func SelectReceiveInto[T any](channelPtr *<-chan T, into *SelectReceiveResult[T]) SelectItem {

	return func(x *selectBuilder) {
		if *channelPtr != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*channelPtr)})
			Append(&x.handlers, func(received reflect.Value, stillOpen bool) bool {
				into.Valid = stillOpen
				into.Value = Cast[T](received, stillOpen)
				into.EndOfStream = !stillOpen
				return false
			})
		}
	}
}

func SelectReceive[T any](channel <-chan T, isValue *bool, value *T, isEndOfStream *bool) SelectItem {
	return func(x *selectBuilder) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
		Append(&x.handlers, func(received reflect.Value, stillOpen bool) bool {
			*isValue = stillOpen
			*value = Cast[T](received, stillOpen)
			*isEndOfStream = !stillOpen
			return false
		})
	}
}

// SelectMustReceive returns done if channel returns endOfStream
func SelectMustReceive[T any](channel <-chan T, isValue *bool, value *T) SelectItem {
	return func(x *selectBuilder) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
		Append(&x.handlers, func(received reflect.Value, stillOpen bool) bool {
			*isValue = stillOpen
			*value = Cast[T](received, stillOpen)
			return !stillOpen
		})
	}
}

func SelectMustReceiveInto[T any](channel <-chan T, into *SelectReceiveResult[T]) SelectItem {
	into.Valid = false
	into.Value = Zero[T]()
	into.EndOfStream = false
	return func(x *selectBuilder) {
		if channel != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
			Append(&x.handlers, func(received reflect.Value, stillOpen bool) bool {
				into.Valid = stillOpen
				into.Value = Cast[T](received, stillOpen)
				into.EndOfStream = !stillOpen
				return !stillOpen
			})
		}
	}
}

func SelectSend[T any](channel chan<- T, value T) SelectItem {
	return func(x *selectBuilder) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channel), Send: reflect.ValueOf(value)})
		Append(&x.handlers, func(_ reflect.Value, _ bool) bool {
			return false
		})
	}
}

func SelectSendWithCallback[T any](channel chan<- T, value T, callback func()) SelectItem {
	return func(x *selectBuilder) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channel), Send: reflect.ValueOf(value)})
		Append(&x.handlers, func(_ reflect.Value, _ bool) bool {
			callback()
			return false
		})
	}
}
