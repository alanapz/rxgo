package rx

import "fmt"

func Take[T any](limit uint) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Void) {

			subscriber, unsubscribe := source.Subscribe()
			defer unsubscribe()

			var position uint

			for {
				msg, isEndOfStream, isDone := Recv1(subscriber, done)

				if isEndOfStream || isDone {
					return
				}

				if !Send1(msg, observer, done) {
					return
				}

				position++

				println(fmt.Sprintf("%d %d", position, limit))

				if position >= limit {
					return
				}
			}
		})
	}
}

func First[T any]() OperatorFunction[T, T] {
	return Take[T](1)
}
