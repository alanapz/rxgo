package rx

import (
	"fmt"
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func TapDebug[T any]() OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(args UnicastObserverArgs[T]) {
			result := drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 source,
				Downstream:             args.Downstream,
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
				NewLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						onSelection: func(msg *u.SelectReceiveMessage[T]) AfterSelectionResult {
							println(fmt.Sprintf("%v value %v", time.Now().Format(time.StampMilli), msg))
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
