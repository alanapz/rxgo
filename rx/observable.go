package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type Observable[T any] interface {
	Subscribe() (<-chan T, func())
}

type AfterSelectionResult string

const ContinueMessage AfterSelectionResult = "continue"
const DropMessage AfterSelectionResult = "drop"
const StopAndContinueNext AfterSelectionResult = "returnContinue"
const StopAndReturnDone AfterSelectionResult = "returnDone"

type drainObservableArgs[T any] struct {
	source         Observable[T]
	downstream     chan<- T
	unsubscribed   <-chan u.Never
	newLoopContext func() drainObservableLoopContext[T]
}

type drainObservableLoopContext[T any] struct {
	beforeSelection func(*[]u.SelectItem) u.SelectResult
	onSelection     func(*u.SelectReceiveMessage[T]) AfterSelectionResult
	beforeSend      func(*u.SelectReceiveMessage[T]) u.SelectResult
	sendValue       func(chan<- T, <-chan u.Never, T) u.SelectResult
}

func drainObservable[T any](args drainObservableArgs[T]) u.SelectResult {

	source, downstream, unsubscribed, newLoopContext := args.source, args.downstream, args.unsubscribed, args.newLoopContext

	u.Require(source, unsubscribed)

	// First subscribe to source observable
	upstream, unsubscribe := source.Subscribe()
	defer unsubscribe()

	for {

		if upstream == nil {
			return u.ContinueResult
		}

		loopCtx := buildNewLoopContext(newLoopContext)

		var msg u.SelectReceiveMessage[T]

		selectionItems := u.Of(u.SelectDone(unsubscribed), u.SelectReceive(&upstream, &msg))

		if loopCtx.beforeSelection(&selectionItems) == u.DoneResult {
			return u.DoneResult
		}

		if u.Selection(selectionItems...) {
			return u.DoneResult
		}

		switch loopCtx.onSelection(&msg) {
		case DropMessage:
			continue
		case StopAndContinueNext:
			return u.ContinueResult
		case StopAndReturnDone:
			return u.DoneResult
		}

		if loopCtx.beforeSend(&msg) == u.DoneResult {
			return u.DoneResult
		}

		if msg.HasValue && loopCtx.sendValue(downstream, unsubscribed, msg.Value) == u.DoneResult {
			return u.DoneResult
		}
	}
}

func buildNewLoopContext[T any](newLoopContext func() drainObservableLoopContext[T]) drainObservableLoopContext[T] {

	loopContext := u.MustCoalesce(newLoopContext, u.Zero[drainObservableLoopContext[T]])()

	return drainObservableLoopContext[T]{
		beforeSelection: u.MustCoalesce(loopContext.beforeSelection, func(_ *[]u.SelectItem) u.SelectResult {
			return u.ContinueResult
		}),
		onSelection: u.MustCoalesce(loopContext.onSelection, func(_ *u.SelectReceiveMessage[T]) AfterSelectionResult {
			return ContinueMessage
		}),
		beforeSend: u.MustCoalesce(loopContext.beforeSend, func(_ *u.SelectReceiveMessage[T]) u.SelectResult {
			return u.ContinueResult
		}),
		sendValue: u.MustCoalesce(loopContext.sendValue, func(downstream chan<- T, unsubscribed <-chan u.Never, value T) u.SelectResult {
			return u.Selection(u.SelectDone(unsubscribed), u.SelectSend(downstream, value))
		}),
	}
}
