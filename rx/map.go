package rx

import u "alanpinder.com/rxgo/v2/utils"

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
	return MapWithError(func(t T) (U, error) {
		return projection(t), nil
	})
}

func MapWithError[T any, U any](projection func(T) (U, error)) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- U, downstreamUnsubscribed <-chan u.Never) error {

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						SendValue: func(_ chan<- T, _ <-chan u.Never, inputValue T) error {

							var outputValue U

							if err := u.Wrap(projection(inputValue))(&outputValue); err != nil {
								return err
							}

							return u.Selection(ctx,
								u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
								u.SelectSend(downstream, outputValue),
							)
						},
					}
				},
			})
		})
	}
}

func MapTo[T any, U any](value U) OperatorFunction[T, U] {
	return Map(func(t T) U { return value })
}
