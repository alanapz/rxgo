package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func FromChannel[T any](channelSupplier func(unsubscribed <-chan u.Never) (upstream <-chan T)) Observable[T] {
	return NewUnicastObservable(func(args UnicastObserverArgs[T]) {

		upstream := channelSupplier(args.DownstreamUnsubscribed)

		for {

			if upstream == nil {
				return
			}

			var msg u.SelectReceiveMessage[T]

			if u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectReceive(&upstream, &msg)) == u.DoneResult {
				return
			}

			if msg.HasValue && u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectSend(args.Downstream, msg.Value)) == u.DoneResult {
				return
			}
		}
	})
}
