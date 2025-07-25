package rx

import (
	"errors"

	u "alanpinder.com/rxgo/v2/utils"
)

type drainObservableArgs[T any] struct {
	Context                *Context
	Source                 Observable[T]
	Downstream             chan<- T
	DownstreamUnsubscribed <-chan u.Never
	NewLoopContext         func() drainObservableLoopContext[T]
}

type drainObservableLoopContext[T any] struct {
	BeforeSelection func(*[]u.SelectItem) error
	AfterSelection  func(*drainObservableMessage[T]) error
	SendValue       func(chan<- T, <-chan u.Never, T) error
}

type drainObservableMessage[T any] struct {
	Valid       bool
	EndOfStream bool
	HasValue    bool
	Value       T
}

var errDropMessage = errors.New("(Drop message)")

func drainObservable[T any](_args drainObservableArgs[T]) error {

	ctx, source, downstream, downstreamUnsubscribed, loopContextBuilder := u.Require(_args.Context), u.Require(_args.Source), _args.Downstream, u.Require(_args.DownstreamUnsubscribed), _args.NewLoopContext

	var upstream <-chan T
	var unsubscribeFromUpstream func()

	if err := u.Wrap2(source.Subscribe(ctx))(&upstream, &unsubscribeFromUpstream); err != nil {
		return err
	}

	defer unsubscribeFromUpstream()

	for {

		if upstream == nil {
			return nil
		}

		loopCtx := buildNewLoopContext(ctx, loopContextBuilder)

		var msg drainObservableMessage[T]

		selectionItems := u.Of(
			u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
			u.SelectReceive(&upstream, func(closed bool, value T) error {
				msg.Valid = true
				msg.EndOfStream = closed
				msg.HasValue = !closed
				msg.Value = value
				return nil
			}),
		)

		if err := loopCtx.BeforeSelection(&selectionItems); err != nil {

			if errors.Is(err, errDropMessage) {
				continue
			}

			if errors.Is(err, ErrEndOfStream) {
				return nil
			}

			return err
		}

		if err := u.Selection(ctx, selectionItems...); err != nil {
			return err
		}

		if err := loopCtx.AfterSelection(&msg); err != nil {

			if errors.Is(err, errDropMessage) {
				continue
			}

			if errors.Is(err, ErrEndOfStream) {
				return nil
			}

			return err
		}

		if msg.HasValue {

			if err := loopCtx.SendValue(downstream, downstreamUnsubscribed, msg.Value); err != nil {
				return err
			}
		}
	}
}

func buildNewLoopContext[T any](ctx *Context, loopContextBuilder func() drainObservableLoopContext[T]) drainObservableLoopContext[T] {

	loopContext := func() drainObservableLoopContext[T] {

		if loopContextBuilder != nil {
			return loopContextBuilder()
		}

		return drainObservableLoopContext[T]{}

	}()

	if loopContext.BeforeSelection == nil {
		loopContext.BeforeSelection = func(_ *[]u.SelectItem) error { return nil }
	}

	if loopContext.AfterSelection == nil {
		loopContext.AfterSelection = func(_ *drainObservableMessage[T]) error { return nil }
	}

	if loopContext.SendValue == nil {

		loopContext.SendValue = func(downstream chan<- T, downstreamUnsubscribed <-chan u.Never, value T) error {
			return u.Selection(ctx, u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)), u.SelectSend(downstream, value))
		}
	}

	return loopContext
}
