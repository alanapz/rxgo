package rx

func Of[T any](values ...T) Observable[T] {
	return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Void) {
		for _, value := range values {
			if !Send1(NewValue(value), observer, done) {
				return
			}
		}
	})
}
