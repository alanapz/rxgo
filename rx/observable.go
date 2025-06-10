package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type Observable[T any] interface {
	Subscribe(env *RxEnvironment) (<-chan T, func(), error)
}

type AfterSelectionResult string

const ContinueMessage AfterSelectionResult = "continue"
const DropMessage AfterSelectionResult = "drop"
const StopAndContinueNext AfterSelectionResult = "returnContinue"
const StopAndReturnDone AfterSelectionResult = "returnDone"

type drainObservableArgs[T any] struct {
	Environment            *RxEnvironment
	Source                 Observable[T]
	Downstream             chan<- T
	DownstreamUnsubscribed <-chan u.Never
	NewLoopContext         func() drainObservableLoopContext[T]
}

type drainObservableLoopContext[T any] struct {
	beforeSelection func(*[]u.SelectItem) u.SelectResult
	onSelection     func(*u.SelectReceiveMessage[T]) AfterSelectionResult
	beforeSend      func(*u.SelectReceiveMessage[T]) u.SelectResult
	sendValue       func(chan<- T, <-chan u.Never, T) u.SelectResult
}

func drainObservable[T any](args drainObservableArgs[T]) u.SelectResult {

	u.Require(args.Environment, args.Source, args.DownstreamUnsubscribed)

	// First subscribe to source observable
	upstream, unsubscribeFromUpstream := args.Source.Subscribe(args.Environment)
	defer unsubscribeFromUpstream()

	for {

		if upstream == nil {
			return u.ContinueResult
		}

		loopCtx := buildNewLoopContext(args.NewLoopContext)

		var msg u.SelectReceiveMessage[T]

		selectionItems := u.Of(u.SelectDone(args.DownstreamUnsubscribed), u.SelectReceive(&upstream, &msg))

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

		if msg.HasValue && loopCtx.sendValue(args.Downstream, args.DownstreamUnsubscribed, msg.Value) == u.DoneResult {
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
		sendValue: u.MustCoalesce(loopContext.sendValue, func(downstream chan<- T, downstreamUnsubscribed <-chan u.Never, value T) u.SelectResult {
			return u.Selection(u.SelectDone(downstreamUnsubscribed), u.SelectSend(downstream, value))
		}),
	}
}
