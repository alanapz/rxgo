package rx

import u "alanpinder.com/rxgo/v2/utils"

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return ConcatMapWithError(func(t T) (Observable[U], error) {
		return projection(t), nil
	})
}

func ConcatMapWithError[T any, U any](projection func(T) (Observable[U], error)) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- U, downstreamUnsubscribed <-chan u.Never) error {

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						SendValue: func(_ chan<- T, _ <-chan u.Never, value T) error {

							var projected Observable[U]

							if err := u.Wrap(projection(value))(&projected); err != nil {
								return err
							}

							return drainObservable(drainObservableArgs[U]{
								Context:                ctx,
								Source:                 projected,
								Downstream:             downstream,
								DownstreamUnsubscribed: downstreamUnsubscribed,
							})
						},
					}
				},
			})
		})
	}
}
