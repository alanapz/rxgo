package rx

import (
	"fmt"
	"maps"
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

type combineLatestCtx[T any] struct {
	valuesOut    chan<- []CombineLatestResult[T]
	errorsOut    chan<- error
	unsubscribed <-chan u.Never
	sources      map[int]*combineLatestSource[T]
}

type combineLatestSource[T any] struct {
	valuesIn    <-chan T
	errorsIn    <-chan error
	unsubscribe func()
	disposed    bool // Whether values AND error are EOF, or done is signalled
	hasValue    bool
	value       T
	complete    bool // Whether values only (not error!) is EOF
}

func CombineLatest[T any](sources ...Observable[T]) Observable[[]CombineLatestResult[T]] {
	return NewUnicastObservable(func(valuesOut chan<- []CombineLatestResult[T], errorsOut chan<- error, unsubscribed <-chan u.Never) {

		ctx := combineLatestCtx[T]{
			valuesOut:    valuesOut,
			errorsOut:    errorsOut,
			unsubscribed: unsubscribed,
			sources:      map[int]*combineLatestSource[T]{},
		}

		for index, source := range sources {
			valuesIn, errorsIn, unsubscribe := source.Subscribe()
			defer unsubscribe()

			ctx.sources[index] = &combineLatestSource[T]{
				valuesIn:    valuesIn,
				errorsIn:    errorsIn,
				unsubscribe: unsubscribe,
			}
		}

		for {

			valueMessages := map[int]*u.SelectReceiveMessage[T]{}
			errorMessages := map[int]*u.SelectReceiveMessage[error]{}

			if ctx.selectNext(valueMessages, errorMessages) == u.DoneResult {
				return
			}

			if ctx.handleErrors(errorMessages) == u.DoneResult {
				return
			}

			if ctx.handleValues(valueMessages) == u.DoneResult {
				return
			}

			if ctx.handleEndOfStream() == u.DoneResult {
				return
			}
		}
	})
}

// combineLatestSelect returns isDone (ie: true for quit, false to keep going)
func (x *combineLatestCtx[T]) selectNext(valueMessages map[int]*u.SelectReceiveMessage[T], errorMessages map[int]*u.SelectReceiveMessage[error]) u.SelectResult {

	selections := u.Of(u.SelectDone(x.unsubscribed))

	for index, source := range x.sources {

		if source.disposed {
			continue
		}

		var valueMsg u.SelectReceiveMessage[T]
		var errorMsg u.SelectReceiveMessage[error]

		if source.valuesIn != nil {
			u.Append(&selections, u.SelectReceive(&source.valuesIn, &valueMsg))
			valueMessages[index] = &valueMsg
		}

		if source.errorsIn != nil {
			u.Append(&selections, u.SelectReceive(&source.errorsIn, &errorMsg))
			errorMessages[index] = &errorMsg
		}
	}

	if len(valueMessages) == 0 && len(errorMessages) == 0 {
		return u.DoneResult
	}

	return u.Selection(selections...)
}

func (x *combineLatestCtx[T]) handleErrors(errorMessages map[int]*u.SelectReceiveMessage[error]) u.SelectResult {

	errors := []error{}

	for msg := range maps.Values(errorMessages) {
		if msg.HasValue {
			u.Append(&errors, msg.Value)
		}
	}

	if len(errors) == 0 {
		return u.ContinueResult
	}

	for err := range slices.Values(errors) {
		if u.Selection(u.SelectDone(x.unsubscribed), u.SelectSend(x.errorsOut, err)) == u.DoneResult {
			return u.DoneResult
		}
	}

	return u.DoneResult
}

func (x *combineLatestCtx[T]) handleValues(valueMessages map[int]*u.SelectReceiveMessage[T]) u.SelectResult {

	var notifyRequired bool

	for index, msg := range valueMessages {
		if msg.HasValue {
			source := x.sources[index]
			source.hasValue = true
			source.value = msg.Value
			notifyRequired = true
		}
	}

	for source := range maps.Values(x.sources) {
		if !source.disposed && source.valuesIn == nil && !source.complete {
			source.complete = true
			notifyRequired = true
		}
	}

	if !notifyRequired {
		return u.ContinueResult
	}

	msgToSend := make([]CombineLatestResult[T], len(x.sources))

	for index, source := range x.sources {
		msgToSend[index] = CombineLatestResult[T]{HasValue: source.hasValue, Value: source.value, Complete: source.complete}
	}

	return u.Selection(u.SelectDone(x.unsubscribed), u.SelectSend(x.valuesOut, msgToSend))
}

func (x *combineLatestCtx[T]) handleEndOfStream() u.SelectResult {

	for source := range maps.Values(x.sources) {
		if !source.disposed && source.valuesIn == nil && source.errorsIn == nil {
			source.disposed = true
			source.unsubscribe()
		}
	}

	for source := range maps.Values(x.sources) {
		if !source.disposed {
			return u.ContinueResult
		}
	}

	return u.DoneResult
}

type CombineLatestResult[T any] struct {
	HasValue bool
	Value    T
	Complete bool
}

func (x CombineLatestResult[T]) String() string {
	if x.HasValue && x.Complete {
		return fmt.Sprintf("{%v, complete}", x.Value)
	}
	if x.HasValue {
		return fmt.Sprintf("{%v}", x.Value)
	}
	if x.Complete {
		return "{complete}"
	}
	return "{}"
}
