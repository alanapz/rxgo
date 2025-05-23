package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Concat[T any](sources ...Observable[T]) Observable[T] {
	return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {
		for source := range slices.Values(sources) {
			if drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
			}) == u.DoneResult {
				return
			}
		}
	})
}
