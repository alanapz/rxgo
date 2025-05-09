package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func TakeUntil[T any, X any](notifiers ...Observable[X]) OperatorFunction[T, T] {

	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(valuesOut chan<- T, errorsOut chan<- error, unsubscribed <-chan u.Never) {

			var notifiersValuesIn []<-chan X

			for notifier := range slices.Values(notifiers) {
				// XXX: We currently ignore notifier errors.. should be propogated or not ?
				notifyValuesIn, _, notifyUnsubscribe := notifier.Subscribe()
				u.Append(&notifiersValuesIn, notifyValuesIn)
				defer notifyUnsubscribe()
			}

			drainObservable(drainObservableArgs[T]{
				source:       source,
				valuesOut:    valuesOut,
				errorsOut:    errorsOut,
				unsubscribed: unsubscribed,
				newLoopContext: func() drainObservableLoopContext[T] {
					return drainObservableLoopContext[T]{
						beforeSelection: func(items *[]u.SelectItem) {
							for notifierValuesIn := range slices.Values(notifiersValuesIn) {
								u.Append(items, u.SelectDone(notifierValuesIn))
							}
						},
					}
				},
			})
		})
	}
}
