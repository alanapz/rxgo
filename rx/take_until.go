package rx

func TakeUntil[T any, X any](notifier <-chan X) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				msg, isEndOfStream, isDone := RecvWithAdditionalDone(subscriber, done, notifier)

				if isEndOfStream || isDone {
					return
				}

				if !Send1(msg, observer, done) {
					return
				}
			}

		})
	}
}
