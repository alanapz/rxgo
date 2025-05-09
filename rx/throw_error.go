package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func ThrowError[T any](err error) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {
		u.Selection(u.SelectDone(unsubscribed), u.SelectSend(errorsOut, err))
	})
}
