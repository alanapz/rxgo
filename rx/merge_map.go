package rx

import (
	"fmt"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func MergeMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(downstream chan<- U, unsubscribed <-chan u.Never) {

			var wg sync.WaitGroup

			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   nil, // Not used
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						sendValue: func(_ chan<- T, _ <-chan u.Never, value T) u.SelectResult {

							wg.Add(1)

							onInnerObservableComplete := u.NewCondition(fmt.Sprintf("Waiting for inner observable for value '%v' to complete", value))

							u.GoRun(func() {
								defer wg.Done()
								defer onInnerObservableComplete()
								drainObservable(drainObservableArgs[U]{
									source:       projection(value),
									downstream:   downstream,
									unsubscribed: unsubscribed,
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
