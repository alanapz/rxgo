package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type Observable[T any] interface {
	Subscribe() (<-chan T, <-chan error, func())
}

type AfterSelectionResult string

const ContinueMessage AfterSelectionResult = "continue"
const DropMessage AfterSelectionResult = "drop"
const StopAndContinueNext AfterSelectionResult = "returnContinue"
const StopAndReturnDone AfterSelectionResult = "returnDone"

type drainObservableArgs[T any] struct {
	source         Observable[T]
	valuesOut      chan<- T
	errorsOut      chan<- error
	unsubscribed   <-chan u.Never
	newLoopContext func() drainObservableLoopContext[T]
}

type drainObservableLoopContext[T any] struct {
	beforeSelection func(*[]u.SelectItem)
	onSelection     func(*u.SelectReceiveMessage[T], *u.SelectReceiveMessage[error]) AfterSelectionResult
	onValue         func(chan<- T, <-chan u.Never, T) u.SelectResult
	onError         func(chan<- error, <-chan u.Never, error) u.SelectResult
}

func drainObservable[T any](args drainObservableArgs[T]) u.SelectResult {

	// First subscribe to source observable
	valuesIn, errorsIn, unsubscribe := args.source.Subscribe()
	defer unsubscribe()

	for {

		if valuesIn == nil && errorsIn == nil {
			return u.ContinueResult
		}

		loopCtx := buildNewLoopContext(args.newLoopContext)

		var valueMsg u.SelectReceiveMessage[T]
		var errorMsg u.SelectReceiveMessage[error]

		selectionItems := u.Of(
			u.SelectDone(args.unsubscribed),
			u.SelectReceive(&valuesIn, &valueMsg),
			u.SelectReceive(&errorsIn, &errorMsg),
		)

		loopCtx.beforeSelection(&selectionItems)

		if u.Selection(selectionItems...) {
			return u.DoneResult
		}

		switch loopCtx.onSelection(&valueMsg, &errorMsg) {
		case DropMessage:
			continue
		case StopAndContinueNext:
			return u.ContinueResult
		case StopAndReturnDone:
			return u.DoneResult
		}

		if valueMsg.HasValue && loopCtx.onValue(args.valuesOut, args.unsubscribed, valueMsg.Value) == u.DoneResult {
			return u.DoneResult
		}

		if errorMsg.HasValue && loopCtx.onError(args.errorsOut, args.unsubscribed, errorMsg.Value) == u.DoneResult {
			return u.DoneResult
		}
	}
}

func buildNewLoopContext[T any](newLoopContext func() drainObservableLoopContext[T]) drainObservableLoopContext[T] {

	loopContext := u.MustCoalesce(newLoopContext, u.Zero[drainObservableLoopContext[T]])()

	return drainObservableLoopContext[T]{
		beforeSelection: u.MustCoalesce(loopContext.beforeSelection, func(_ *[]u.SelectItem) {}),
		onSelection: u.MustCoalesce(loopContext.onSelection, func(_ *u.SelectReceiveMessage[T], _ *u.SelectReceiveMessage[error]) AfterSelectionResult {
			return ContinueMessage
		}),
		onValue: u.MustCoalesce(loopContext.onValue, func(valuesOut chan<- T, unsubscribed <-chan u.Never, value T) u.SelectResult {
			return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, value))
		}),
		onError: u.MustCoalesce(loopContext.onError, func(errorsOut chan<- error, unsubscribed <-chan u.Never, err error) u.SelectResult {
			return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(errorsOut, err))
		}),
	}
}
