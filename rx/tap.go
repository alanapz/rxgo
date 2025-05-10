package rx

import u "alanpinder.com/rxgo/v2/utils"

func Tap[T any](tapValue func(T)) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {
			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {

							if tapValue != nil && msg.HasValue {
								tapValue(msg.Value)
							}

							return ContinueMessage
						},
					}
				},
			})
		})
	}
}
