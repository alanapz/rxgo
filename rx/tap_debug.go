package rx

import (
	"fmt"
)

func TapDebug[T any](label string, writer func(message ...any) (int, error)) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				writer(fmt.Sprintf("%s -> Waiting for next message ...", label))

				msg, isEndOfStream, isDone := Recv1(subscriber, done)

				if isEndOfStream || isDone {
					writer(fmt.Sprintf("%s -> Closing as source observable is EOF", label))
					return
				}

				if isDone {
					writer(fmt.Sprintf("%s -> Closing as downstream subscriber has signalled done", label))
					return
				}

				writer(fmt.Sprintf("%s -> Received %s, sending to observer ...", label, msg))

				if !Send1(msg, observer, done) {
					writer(fmt.Sprintf("%s -> Closing as downstream observer / subscriber has signalled done (or closed stream)", label))
					return
				}

				writer(fmt.Sprintf("%s -> Sent to observer: %s", label, msg))
			}
		})
	}
}
