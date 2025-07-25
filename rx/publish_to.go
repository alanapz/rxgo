package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func PublishTo[T any](ctx *Context, source Observable[T], sink Subject[T], endOfStreamPropagation EndOfStreamPropagationPolicy) error {

	var sinkEndOfStream chan u.Never
	var sendSinkEndOfStream func()

	if err := u.Wrap2(NewChannel[u.Never](ctx, 0))(&sinkEndOfStream, &sendSinkEndOfStream); err != nil {
		return err
	}

	GoRun(ctx, func() {

		defer sendSinkEndOfStream()

		// Important: Sink can become end-of-stream via other publishers
		// We need to thus add a listener for sink end-of-strea
		defer sink.OnEndOfStream(sendSinkEndOfStream)()

		drainObservable(drainObservableArgs[T]{
			Context:                ctx,
			Source:                 source,
			Downstream:             nil, // Not used
			DownstreamUnsubscribed: sinkEndOfStream,
			NewLoopContext: func() drainObservableLoopContext[T] {
				return drainObservableLoopContext[T]{

					AfterSelection: func(msg *drainObservableMessage[T]) error {

						if msg.Valid && msg.HasValue {

							if err := sink.Next(msg.Value); err != nil {
								return err
							}
						}

						if msg.Valid && msg.EndOfStream && endOfStreamPropagation == PropogateEndOfStream {

							if err := sink.EndOfStream(); err != nil {
								return err
							}
						}

						return errDropMessage
					},
				}
			},
		})
	})

	return nil
}
