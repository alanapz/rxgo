package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func FromChannel[T any](channelSupplier func(unsubscribed <-chan u.Never) (upstream <-chan T)) Observable[T] {
	return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {

		upstream := channelSupplier(unsubscribed)

		for {

			if upstream == nil {
				return
			}

			var msg u.SelectReceiveMessage[T]

			if u.Selection(u.SelectDone(unsubscribed), u.SelectReceive(&upstream, &msg)) == u.DoneResult {
				return
			}

			if msg.HasValue && u.Selection(u.SelectDone(unsubscribed), u.SelectSend(downstream, msg.Value)) == u.DoneResult {
				return
			}
		}
	})
}
