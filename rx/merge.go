package rx

import (
	"fmt"

	u "alanpinder.com/rxgo/v2/utils"
)

func Merge[T any](sources ...Observable[T]) Observable[T] {

	if len(sources) == 0 {
		return Of[T]()
	}

	if len(sources) == 1 {
		return sources[0]
	}

	return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

		var wg u.WaitGroup

		wg.Add(len(sources))

		for index, source := range sources {

			label := fmt.Sprintf("Waiting for inner observable for source #%d to complete", index)

			GoRunWithLabel(ctx, label, func() {

				defer wg.Done()

				err := drainObservable(drainObservableArgs[T]{
					Context:                ctx,
					Source:                 source,
					Downstream:             downstream,
					DownstreamUnsubscribed: downstreamUnsubscribed,
				})

				if err != nil {
					wg.Error(err)
				}
			})
		}

		wg.Wait()

		return wg.GetError()
	})
}
