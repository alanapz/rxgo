package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Of[T any](values ...T) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {
		for value := range slices.Values(values) {
			if u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, value)) == u.DoneResult {
				return
			}
		}
	})
}
