package rx

import u "alanpinder.com/rxgo/v2/utils"

func Filter[T any](filter func(T) bool) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

							if msg.HasValue && !filter(msg.Value) {
								return DropMessage
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}
