package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func PublishTo[T any](env *RxEnvironment, source Observable[T], sink Subject[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {

	u.Require(env, source, sink)

	sinkEndofStream, closeSinkEndOfStreamChannel := NewChannel[u.Never](env, 0)

	u.GoRun(func() {

		defer closeSinkEndOfStreamChannel.Emit()

		// Important: Sink can become end-of-stream via other publishers
		// We need to thus add a listener for sink end-of-strea
		defer sink.OnEndOfStream(closeSinkEndOfStreamChannel.Emit)()

		drainObservable(drainObservableArgs[T]{
			Environment:            env,
			Source:                 source,
			Downstream:             nil, // Not used
			DownstreamUnsubscribed: sinkEndofStream,
			NewLoopContext: func() drainObservableLoopContext[T] {
				return drainObservableLoopContext[T]{
					onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

						if msg.Selected && msg.HasValue {
							env.Error(sink.Next(msg.Value))
						}

						if msg.Selected && msg.EndOfStream && endOfStreamPropagation == PropogateEndOfStream {
							env.Error(sink.EndOfStream())
						}

						return DropMessage
					},
				}
			},
		})
	})
}
