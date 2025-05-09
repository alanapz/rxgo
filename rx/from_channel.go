package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func FromChannel[T any](channelSupplier func(unsubscribed <-chan u.Never) (valuesIn <-chan T, errorsIn <-chan error)) Observable[T] {
	return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {

		valuesIn, errorsIn := channelSupplier(unsubscribed)

		for {

			if valuesIn == nil && errorsIn == nil {
				return
			}

			var valueMsg u.SelectReceiveMessage[T]
			var errorMsg u.SelectReceiveMessage[error]

			if u.Selection(u.SelectDone(unsubscribed), u.SelectReceive(&valuesIn, &valueMsg), u.SelectReceive(&errorsIn, &errorMsg)) == u.DoneResult {
				return
			}

			if valueMsg.HasValue && u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, valueMsg.Value)) == u.DoneResult {
				return
			}

			if errorMsg.HasValue && u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, valueMsg.Value)) == u.DoneResult {
				return
			}
		}
	})
}
