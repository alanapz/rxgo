package rx

func Map[T any, U any](projection func(T) U) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(observer chan<- Message[U], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {
				msg, isEndOfStream, isDone := Recv1(subscriber, done)

				if isEndOfStream || isDone {
					return
				}

				if !Send1(MapMessage(msg, projection), observer, done) {
					return
				}
			}
		})
	}
}

func ToAny[T any]() OperatorFunction[T, any] {
	return Map(func(t T) any { return t })
}

func FromAny[T any]() OperatorFunction[any, T] {
	return Map(func(t any) T { return t.(T) })
}
