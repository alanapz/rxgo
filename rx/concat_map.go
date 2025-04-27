package rx

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(observer chan<- Message[U], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				msg, isEndOfStream, isDone := Recv1(subscriber, done)

				if isEndOfStream || isDone {
					return
				}

				if !drainObservable(projection(msg.Value), observer, done) {
					return
				}
			}
		})
	}
}

func drainObservable[U any](observable Observable[U], observer chan<- Message[U], done <-chan Void) bool {

	subscription, unsubscribe := observable.Subscribe()
	defer unsubscribe()

	for {

		msg, isEndOfStream, isDone := Recv1(subscription, done)

		if isEndOfStream {
			return true
		}

		if isDone {
			return false
		}

		if !Send1(msg, observer, done) {
			// If returns false here, means upstream observer has unsubscribed
			return false
		}
	}
}
