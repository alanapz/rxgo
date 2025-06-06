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

	return NewUnicastObservable(func(args UnicastObserverArgs[T]) {
		for source := range slices.Values(sources) {
			if drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             args.Downstream,
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
			}) == u.DoneResult {
				return
			}
		}
	})
}
