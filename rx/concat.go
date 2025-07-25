package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Concat[T any](sources ...Observable[T]) Observable[T] {

	if len(sources) == 0 {
		return Of[T]()
	}

	if len(sources) == 1 {
		return sources[0]
	}

	return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

		for source := range slices.Values(sources) {

			err := drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 source,
				Downstream:             downstream,
				DownstreamUnsubscribed: downstreamUnsubscribed,
			})

			if err != nil {
				return err
			}
		}

		return nil
	})
}
