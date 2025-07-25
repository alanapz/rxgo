package rx

import u "alanpinder.com/rxgo/v2/utils"

func Tap[T any](tapConsumer func(T) error) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             downstream,
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						AfterSelection: func(msg *drainObservableMessage[T]) error {

							if tapConsumer != nil && msg.HasValue {

								if err := tapConsumer(msg.Value); err != nil {
									return err
								}

							}

							return nil
						},
					}
				},
			})
		})
	}
}
