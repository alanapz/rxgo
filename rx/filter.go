package rx

import u "alanpinder.com/rxgo/v2/utils"

func Filter[T any](filter func(T) (bool, error)) OperatorFunction[T, T] {
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

							if !msg.HasValue {
								return nil
							}

							var keepMessage bool

							if err := u.Wrap(filter(msg.Value))(&keepMessage); err != nil {
								return err
							}

							if !keepMessage {
								return errDropMessage
							}

							return nil
						},
					}
				},
			})
		})
	}
}
