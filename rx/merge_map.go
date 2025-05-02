package rx

import (
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func MergeMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, done <-chan u.Never) {

			var wg sync.WaitGroup

			drainObservable(drainObservableArgs[T]{
				source:    source,
				valuesOut: nil, // Not used
				errorsOut: errorsOut,
				done:      done,
				valueHandler: func(_ chan<- T, done <-chan u.Never, value T) u.SelectResult {

					wg.Add(1)

					GoRun(func() {
						defer wg.Done()
						drainObservable(drainObservableArgs[U]{source: projection(value), valuesOut: valuesOut, errorsOut: errorsOut, done: done})
					})

					return u.ContinueResult // Cant do better, merge is async
				},
			})

			wg.Wait()
		})
	}
}
