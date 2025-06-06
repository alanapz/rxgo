package rx

import (
	"fmt"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func Merge[T any](sources ...Observable[T]) Observable[T] {

	if len(sources) == 0 {
		return Of[T]()
	}

	if len(sources) == 1 {
		return sources[0]
	}

	return NewUnicastObservable(func(args UnicastObserverArgs[T]) {

		var wg sync.WaitGroup

		wg.Add(len(sources))

		for index, source := range sources {

			onInnerObservableComplete := u.NewCondition(fmt.Sprintf("Waiting for inner observable for source #%d to complete", index))

			args.Environment.Execute(func() {
				defer wg.Done()
				defer onInnerObservableComplete()
				drainObservable(drainObservableArgs[T]{
					Environment:            args.Environment,
					Source:                 source,
					Downstream:             args.Downstream,
					DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				})
			})
		}

		wg.Wait()
	})
}
