package rx

// func FromUnmanagedChannel[T any](source <-chan T) Observable[T] {
// 	return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Never) {
// 		for {

// 			value, sourceClosed, isDone := Recv1(source, done)

// 			if sourceClosed || isDone {
// 				return
// 			}

// 			if !Send1(NewValue(value), observer, done) {
// 				return
// 			}
// 		}
// 	})
// }

// func FromChannel[T any](sourceSupplier func() (<-chan T, func())) Observable[T] {
// 	return NewUnicastObservable(func(observer chan<- Message[T], done <-chan Never) {

// 		source, cleanupSource := sourceSupplier()
// 		defer cleanupSource()

// 		for {

// 			value, sourceClosed, isDone := Recv1(source, done)

// 			if sourceClosed || isDone {
// 				return
// 			}

// 			if !Send1(NewValue(value), observer, done) {
// 				return
// 			}
// 		}
// 	})
// }
