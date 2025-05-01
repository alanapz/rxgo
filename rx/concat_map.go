package rx

func ConcatMap[T any, U any](projection func(T) Observable[U]) OperatorFunction[T, U] {
	return func(source Observable[T]) Observable[U] {
		return NewUnicastObservable(func(valuesOut chan<- U, errorsOut chan<- error, done <-chan Never) {

			valuesIn, errorsIn, unsubscribe := source.Subscribe()
			defer unsubscribe()

			for {

				if valuesIn == nil && errorsIn == nil {
					return
				}

				var valueMsg SelectReceiveResult[T]
				var errorMsg SelectReceiveResult[error]

				if Selection(SelectDone(done), SelectReceiveInto(valuesIn, &valueMsg), SelectReceiveInto(errorsIn, &errorMsg)) {
					// If Selection returns true, means done was signalled, so return true (done signalled)
					return
				}

				if valueMsg.EndOfStream {
					valuesIn = nil
				}

				if errorMsg.EndOfStream {
					errorsIn = nil
				}

				if valueMsg.Valid && drainObservable(projection(valueMsg.Value), valuesOut, errorsOut, done) {
					return
				}

				if errorMsg.Valid && Selection(SelectDone(done), SelectSend(errorsOut, errorMsg.Value)) {
					return
				}
			}
		})
	}
}

// drainObservable returns true if done was signalled, false if source (both values and error) was EOF
func drainObservable[U any](source Observable[U], valuesOut chan<- U, errorsOut chan<- error, done <-chan Never) bool {

	valuesIn, errorsIn, unsubscribe := source.Subscribe()
	defer unsubscribe()

	for {

		if valuesIn == nil && errorsIn == nil {
			println("drainObservable closed valuesIn errorsIn nil")
			return false
		}

		var valueMsg SelectReceiveResult[U]
		var errorMsg SelectReceiveResult[error]

		if Selection(SelectDone(done), SelectReceiveInto(valuesIn, &valueMsg), SelectReceiveInto(errorsIn, &errorMsg)) {
			// If Selection returns true, means done was signalled, so return true (done signalled)
			println("drainObservable closed")
			return true
		}

		if valueMsg.EndOfStream {
			println("drainObservable valueMsg.EndOfStream")
			valuesIn = nil
		}

		if errorMsg.EndOfStream {
			println("drainObservable errorMsg.EndOfStream")
			errorsIn = nil
		}

		if valueMsg.Valid && Selection(SelectDone(done), SelectSend(valuesOut, valueMsg.Value)) {
			println("drainObservable closed valueMsg.Valid ")
			return true
		}

		if errorMsg.Valid && Selection(SelectDone(done), SelectSend(errorsOut, errorMsg.Value)) {
			println("drainObservable closed errorMsg.Valid")
			return true
		}
	}
}
