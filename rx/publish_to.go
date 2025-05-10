package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type PublishToArgs[T any] struct {
	Source               Observable[T]
	Sink                 Subject[T]
	PropogateEndOfStream bool
}

func PublishTo[T any](args PublishToArgs[T]) {

	source, sink, propogateEndOfStream := args.Source, args.Sink, args.PropogateEndOfStream

	u.Require(source, sink)

	var unsubscribedCleanup u.Event

	unsubscribed := u.NewChannel[u.Never](&unsubscribedCleanup, 0)

	u.GoRun(func() {

		defer unsubscribedCleanup.Emit()
		defer sink.OnEndOfStream(unsubscribedCleanup.Emit)()

		drainObservable(drainObservableArgs[T]{
			source:       source,
			downstream:   nil, // Not used
			unsubscribed: unsubscribed,
			newLoopContext: func() drainObservableLoopContext[T] {
				return drainObservableLoopContext[T]{
					onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

						if msg.Selected && msg.HasValue {
							sink.Next(msg.Value)
						}

						if msg.Selected && msg.EndOfStream && propogateEndOfStream {
							sink.EndOfStream()
						}

						return DropMessage
					},
				}
			},
		})
	})
}
