package utils

import (
	"reflect"
	"slices"
)

type SelectResult bool

const ContinueResult SelectResult = false
const DoneResult SelectResult = true

type selectBuilder struct {
	cases    []reflect.SelectCase
	handlers []func(reflect.Value, bool) SelectResult
}

type SelectItem func(*selectBuilder) SelectResult

// Selection returns done (or end of stream)
func Selection(selections ...SelectItem) SelectResult {

	builder := &selectBuilder{
		cases:    []reflect.SelectCase{},
		handlers: []func(reflect.Value, bool) SelectResult{},
	}

	for selection := range slices.Values(selections) {
		if result := selection(builder); result == DoneResult {
			return DoneResult
		}
	}

	chosen, recv, open := reflect.Select(builder.cases)
	return builder.handlers[chosen](recv, open)
}

// SelectDone returns endOfStream if channel is not open
func SelectDone[T any](channel <-chan T) SelectItem {
	return func(x *selectBuilder) SelectResult {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(channel)})
		Append(&x.handlers, func(_ reflect.Value, stillOpen bool) SelectResult {
			return Ternary(stillOpen, ContinueResult, DoneResult)
		})
		return ContinueResult
	}
}

func SelectSend[T any](channel chan<- T, value T) SelectItem {
	return func(x *selectBuilder) SelectResult {
		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(channel), Send: reflect.ValueOf(value)})
		Append(&x.handlers, func(_ reflect.Value, _ bool) SelectResult {
			return ContinueResult
		})
		return ContinueResult
	}
}
