package rx

import (
	"reflect"
)

type MessageSelectResult[T any] struct {
	Quit bool
	Msg  Message[T]
}

type Select struct {
	cases    []reflect.SelectCase
	handlers []func(reflect.Value, bool)
}

func NewSelect() *Select {
	return &Select{
		cases:    []reflect.SelectCase{},
		handlers: []func(reflect.Value, bool){},
	}
}

func (x *Select) Add(item func(*Select)) {
	item(x)
}

func (x *Select) Select() {
	chosen, recv, open := reflect.Select(x.cases)
	x.handlers[chosen](recv, open)
}

func RevcItem[T any](channel <-chan T, callback func(T, bool)) func(*Select) {
	return func(x *Select) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
		Append(&x.handlers, func(value reflect.Value, stillOpen bool) {
			callback(value.Interface().(T), !stillOpen)
		})
	}
}

func SendItem[T any](channel chan<- T, value T, callback func()) func(*Select) {
	return func(x *Select) {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channel), Send: reflect.ValueOf(value)})
		Append(&x.handlers, func(_ reflect.Value, _ bool) {
			callback()
		})
	}
}
