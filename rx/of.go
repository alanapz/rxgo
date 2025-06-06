package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Of[T any](values ...T) Observable[T] {
	return NewUnicastObservable(func(args UnicastObserverArgs[T]) {
		for value := range slices.Values(values) {
			if u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectSend(args.Downstream, value)) == u.DoneResult {
				return
			}
		}
	})
}
