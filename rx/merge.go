package rx

import (
	"slices"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func Merge[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never) {

		var wg sync.WaitGroup

		for source := range slices.Values(sources) {

			wg.Add(1)

			GoRun(func() {
				defer wg.Done()
				drainObservable(drainObservableArgs[T]{source: source, valuesOut: valuesOut, errorsOut: errorsOut, done: done})
			})
		}

		wg.Wait()
	})
}
