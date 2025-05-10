package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func TakeUntil[T any, X any](notifiers ...Observable[X]) OperatorFunction[T, T] {

	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(downstream chan<- T, unsubscribed <-chan u.Never) {

			var notifierSources []<-chan X

			for notifier := range slices.Values(notifiers) {
				notifierSource, notifierUnsubscribe := notifier.Subscribe()
				u.Append(&notifierSources, notifierSource)
				defer notifierUnsubscribe()
			}

			drainObservable(drainObservableArgs[T]{
				source:       source,
				downstream:   downstream,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						beforeSelection: func(items *[]u.SelectItem) u.SelectResult {
							for notifierSource := range slices.Values(notifierSources) {
								u.Append(items, u.SelectDone(notifierSource))
							}
							return u.ContinueResult
						},
					}
				},
			})
		})
	}
}
