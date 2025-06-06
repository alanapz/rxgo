package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type PublishToArgs[T any] struct {
	Environment          *RxEnvironment
	Source               Observable[T]
	Sink                 Subject[T]
	PropogateEndOfStream bool
}

func PublishTo[T any](args PublishToArgs[T]) {

	u.Require(args.Environment, args.Source, args.Sink)

	sinkEndofStream, closeSinkEndOfStreamChannel := NewChannel[u.Never](args.Environment, 0)

	u.GoRun(func() {

		defer closeSinkEndOfStreamChannel.Resolve()

		// Important: Sink can become end-of-stream via other publishers
		// We need to thus add a listener for sink end-of-strea
		defer args.Sink.OnEndOfStream(closeSinkEndOfStreamChannel.Resolve)()

		drainObservable(drainObservableArgs[T]{
			Environment:            args.Environment,
			Source:                 args.Source,
			Downstream:             nil, // Not used
			DownstreamUnsubscribed: sinkEndofStream,
			NewLoopContext: func() drainObservableLoopContext[T] {
				return drainObservableLoopContext[T]{
					onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

						if msg.Selected && msg.HasValue {
							args.Environment.Error(args.Sink.Next(msg.Value))
						}

						if msg.Selected && msg.EndOfStream && args.PropogateEndOfStream {
							args.Environment.Error(args.Sink.EndOfStream())
						}

						return DropMessage
					},
				}
			},
		})
	})
}
