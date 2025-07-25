package rx

import (
	"fmt"

	u "alanpinder.com/rxgo/v2/utils"
)

func MergeMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return MergeMapWithError(func(t T) (Observable[U], error) {
		return projection(t), nil
	})
}

func MergeMapWithError[T any, U any](projection func(T) (Observable[U], error)) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- U, downstreamUnsubscribed <-chan u.Never) error {

			var wg u.WaitGroup

			err := drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: downstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {

					return drainObservableLoopContext[T]{

						SendValue: func(_ chan<- T, _ <-chan u.Never, value T) error {

							wg.Add(1)

							label := fmt.Sprintf("Waiting for inner observable for value '%v' to complete", value)

							GoRunWithLabel(ctx, label, func() {

								defer wg.Done()

								var projected Observable[U]

								if err := u.Wrap(projection(value))(&projected); err != nil {
									wg.Error(fmt.Errorf("projecting '%v' to merge-map observable: %w", value, err))
									return
								}

								err := drainObservable(drainObservableArgs[U]{
									Context:                ctx,
									Source:                 projected,
									Downstream:             downstream,
									DownstreamUnsubscribed: downstreamUnsubscribed,
								})

								if err != nil {
									wg.Error(fmt.Errorf("draining merge-map projection for '%v': %w", value, err))
								}
							})

							return nil // Cant do better, merge is async
						},
					}
				},
			})

			wg.Wait()

			if err != nil {
				wg.Error(err)
			}

			return wg.GetError()
		})
	}
}
