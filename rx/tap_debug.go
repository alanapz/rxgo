package rx

import (
	"fmt"

	u "alanpinder.com/rxgo/v2/utils"
)

func TapDebug[T any]() OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {
			result := drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {
							println(fmt.Sprintf("value %v", msg))
							return ContinueMessage
						},
					}
				},
			})

			if result == u.ContinueResult {
				println("quit by end of stream")
			}
			if result == u.DoneResult {
				println("quit by unsubscribe")
			}
		})
	}
}
