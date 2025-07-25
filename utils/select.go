package utils

import (
	"errors"
	"reflect"
	"slices"
)

var ErrDisposed = errors.New("disposed")
var ErrDisposalInProgress = errors.New("disposal in progress")

type cancellableContext interface {
	Done() <-chan struct{}
	Err() error
}

type selectBuilder struct {
	cases    []reflect.SelectCase
	handlers []func(reflect.Value, bool) error
}

type SelectItem func(*selectBuilder) error

// Selection returns done (or end of stream)
// XXXX: Maybe remove ctx? as cant control error returned
func Selection(ctx cancellableContext, selections ...SelectItem) error {

	builder := &selectBuilder{
		cases:    []reflect.SelectCase{},
		handlers: []func(reflect.Value, bool) error{},
	}

	SelectDone(ctx.Done(), ctx.Err)

	for selection := range slices.Values(selections) {
		if err := selection(builder); err != nil {
			return err
		}
	}

	chosen, recv, open := reflect.Select(builder.cases)
	return builder.handlers[chosen](recv, open)
}

// Selection returns done (or end of stream)
// XXXX: Maybe remove ctx? as cant control error returned
func Selection2(selections ...SelectItem) error {

	builder := &selectBuilder{
		cases:    []reflect.SelectCase{},
		handlers: []func(reflect.Value, bool) error{},
	}

	for selection := range slices.Values(selections) {
		if err := selection(builder); err != nil {
			return err
		}
	}

	chosen, recv, open := reflect.Select(builder.cases)
	return builder.handlers[chosen](recv, open)
}

// SelectDisposed returns ErrDisposed if channel is not open or becomes signalled
func SelectDisposed(channel <-chan Void) SelectItem {
	return func(x *selectBuilder) error {

		if channel != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
			Append(&x.handlers, func(_ reflect.Value, stillOpen bool) error { return ErrDisposed })
		}

		return nil
	}
}

// SelectDone returns endOfStream if channel is not open or becomes signalled
func SelectDone[T any](channel <-chan T, errSupplier func() error) SelectItem {
	return func(x *selectBuilder) error {

		if channel != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
			Append(&x.handlers, func(_ reflect.Value, stillOpen bool) error { return errSupplier() })
		}

		return nil
	}
}

func SelectSend[T any](channel chan<- T, value T) SelectItem {
	return func(x *selectBuilder) error {

		if channel != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channel), Send: reflect.ValueOf(value)})
			Append(&x.handlers, func(_ reflect.Value, _ bool) error { return nil })
		}

		return nil
	}
}

func SelectReceive[T any](channel *<-chan T, handler func(closed bool, value T) error) SelectItem {
	return func(x *selectBuilder) error {

		if *channel != nil {
			Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*channel)})
			Append(&x.handlers, func(received reflect.Value, stillOpen bool) error {

				if !stillOpen {
					*channel = nil
					return handler(true, Zero[T]())
				}

				return handler(false, received.Interface().(T))
			})
		}

		return nil
	}
}
