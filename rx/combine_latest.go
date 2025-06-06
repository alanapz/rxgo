package rx

import (
	"fmt"
	"maps"

	u "alanpinder.com/rxgo/v2/utils"
)

type combineLatestCtx[T any] struct {
	downstream             chan<- []CombineLatestResult[T]
	downstreamUnsubscribed <-chan u.Never
	sources                map[int]*combineLatestSource[T]
}

type combineLatestSource[T any] struct {
	upstream    <-chan T
	unsubscribe func()
	disposed    bool
	hasValue    bool
	value       T
	endOfStream bool
}

func CombineLatest[T any](sources ...Observable[T]) Observable[[]CombineLatestResult[T]] {
	return NewUnicastObservable(func(args UnicastObserverArgs[[]CombineLatestResult[T]]) {

		ctx := combineLatestCtx[T]{
			downstream:             args.Downstream,
			downstreamUnsubscribed: args.DownstreamUnsubscribed,
			sources:                map[int]*combineLatestSource[T]{},
		}

		for index, source := range sources {
			upstream, unsubscribe := source.Subscribe(args.Environment)
			defer unsubscribe()

			ctx.sources[index] = &combineLatestSource[T]{
				upstream:    upstream,
				unsubscribe: unsubscribe,
			}
		}

		for {

			messages := map[int]*u.SelectReceiveMessage[T]{}

			if ctx.selectNext(messages) == u.DoneResult {
				return
			}

			if ctx.handleMessages(messages) == u.DoneResult {
				return
			}

			if ctx.handleEndOfStream() == u.DoneResult {
				return
			}
		}
	})
}

// combineLatestSelect returns isDone (ie: true for quit, false to keep going)
func (x *combineLatestCtx[T]) selectNext(messages map[int]*u.SelectReceiveMessage[T]) u.SelectResult {

	selections := u.Of(u.SelectDone(x.downstreamUnsubscribed))

	for index, source := range x.sources {

		if source.disposed {
			continue
		}

		var msg u.SelectReceiveMessage[T]

		if source.upstream != nil {
			u.Append(&selections, u.SelectReceive(&source.upstream, &msg))
			messages[index] = &msg
		}
	}

	if len(messages) == 0 {
		return u.DoneResult
	}

	return u.Selection(selections...)
}

func (x *combineLatestCtx[T]) handleMessages(messages map[int]*u.SelectReceiveMessage[T]) u.SelectResult {

	var notifyRequired bool

	for index, msg := range messages {
		if msg.HasValue {
			source := x.sources[index]
			source.hasValue = true
			source.value = msg.Value
			notifyRequired = true
		}
	}

	for source := range maps.Values(x.sources) {
		if source.upstream == nil && !source.endOfStream {
			source.endOfStream = true
			notifyRequired = true
		}
	}

	if !notifyRequired {
		return u.ContinueResult
	}

	msgToSend := make([]CombineLatestResult[T], len(x.sources))

	for index, source := range x.sources {
		msgToSend[index] = CombineLatestResult[T]{HasValue: source.hasValue, Value: source.value, EndOfStream: source.endOfStream}
	}

	return u.Selection(u.SelectDone(x.downstreamUnsubscribed), u.SelectSend(x.downstream, msgToSend))
}

func (x *combineLatestCtx[T]) handleEndOfStream() u.SelectResult {

	for source := range maps.Values(x.sources) {
		if !source.disposed && source.endOfStream {
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
	HasValue    bool
	Value       T
	EndOfStream bool
}

func (x CombineLatestResult[T]) String() string {
	if x.HasValue && x.EndOfStream {
		return fmt.Sprintf("{%v, EOF}", x.Value)
	}
	if x.HasValue {
		return fmt.Sprintf("{%v}", x.Value)
	}
	if x.EndOfStream {
		return "{EOF}"
	}
	return "{}"
}
