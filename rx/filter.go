package rx

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				msg, isEndOfStream, isDone := Recv1(subscriber, done)

				if isEndOfStream || isDone {
					return
				}

				if msg.IsValue() && !filter(msg.Value) {
					continue
				}

				if !Send1(msg, observer, done) {
					return
				}
			}
		})
	}
}
