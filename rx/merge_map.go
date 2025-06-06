package rx

import (
	"fmt"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func MergeMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(args UnicastObserverArgs[U]) {

			var wg sync.WaitGroup

			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             nil, // Not used
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {

							wg.Add(1)

							onInnerObservableComplete := u.NewCondition(fmt.Sprintf("Waiting for inner observable for value '%v' to complete", value))

							args.Environment.Execute(func() {
								defer wg.Done()
								defer onInnerObservableComplete()
								drainObservable(drainObservableArgs[U]{
									Environment:            args.Environment,
									Source:                 projection(value),
									Downstream:             args.Downstream,
									DownstreamUnsubscribed: args.DownstreamUnsubscribed,
								})
							})

							return u.ContinueResult // Cant do better, merge is async
						},
					}
				},
			})

			wg.Wait()
		})
	}
}
