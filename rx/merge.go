package rx

import (
	"fmt"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func Merge[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {

		var wg sync.WaitGroup

		for index, source := range sources {

			wg.Add(1)

			onInnerObservableComplete := u.NewCondition(fmt.Sprintf("Waiting for inner observable for source #%d to complete", index))

			u.GoRun(func() {
				defer wg.Done()
				defer onInnerObservableComplete()
				drainObservable(drainObservableArgs[T]{
					source:       source,
					downstream:   downstream,
					unsubscribed: unsubscribed,
				})
			})
		}

		wg.Wait()
	})
}
