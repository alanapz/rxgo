package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

type Observable[T any] interface {
	Subscribe() (<-chan T, <-chan error, func())
}

type AfterSelectionResult string

const ContinueMessage AfterSelectionResult = "continue"
const DropMessage AfterSelectionResult = "drop"
const ReturnContinue AfterSelectionResult = "returnContinue"
const ReturnDone AfterSelectionResult = "returnDone"

type drainObservableArgs[T any] struct {
	source                Observable[T]
	valuesOut             chan<- T
	errorsOut             chan<- error
	done                  <-chan u.Never
	additionalSelections  func() []u.SelectItem
	afterSelectionHandler func(*u.SelectReceiveMessage[T], *u.SelectReceiveMessage[error]) AfterSelectionResult
	valueHandler          func(chan<- T, <-chan u.Never, T) u.SelectResult
	errorHandler          func(chan<- error, <-chan u.Never, error) u.SelectResult
}

func drainObservable[T any](args drainObservableArgs[T]) u.SelectResult {

	source, valuesOut, errorsOut, done, additionalSelections, afterSelection, valueHandler, errorHandler :=
		args.source,
		args.valuesOut,
		args.errorsOut,
		args.done,
		u.MustCoalesce(args.additionalSelections, defaultAdditionalSelections),
		u.MustCoalesce(args.afterSelectionHandler, defaultAfterSelectionHandler),
		u.MustCoalesce(args.valueHandler, defaultDrainValueHandler),
		u.MustCoalesce(args.errorHandler, defaultDrainErrorHandler)

	valuesIn, errorsIn, unsubscribe := source.Subscribe()
	defer unsubscribe()

	for {

		if valuesIn == nil && errorsIn == nil {
			return u.ContinueResult
		}

		var valueMsg u.SelectReceiveMessage[T]
		var errorMsg u.SelectReceiveMessage[error]

		selectionItems := u.Of(
			u.SelectDone(done),
			u.SelectReceive(&valuesIn, &valueMsg),
			u.SelectReceive(&errorsIn, &errorMsg),
		)

		u.Append(&selectionItems, additionalSelections()...)

		if u.Selection(selectionItems...) {
			return u.DoneResult
		}

		switch afterSelection(&valueMsg, &errorMsg) {
		case DropMessage:
			continue
		case ReturnContinue:
			return u.ContinueResult
		case ReturnDone:
			return u.DoneResult
		}

		if valueMsg.HasValue && valueHandler(valuesOut, done, valueMsg.Value) == u.DoneResult {
			return u.DoneResult
		}

		if errorMsg.HasValue && errorHandler(errorsOut, done, errorMsg.Value) == u.DoneResult {
			return u.DoneResult
		}
	}
}

func defaultAdditionalSelections() []u.SelectItem {
	return u.Zero[[]u.SelectItem]()
}

func defaultAfterSelectionHandler[T any](valueMsg *u.SelectReceiveMessage[T], errorMsg *u.SelectReceiveMessage[error]) AfterSelectionResult {
	return ContinueMessage
}

func defaultDrainValueHandler[T any](valuesOut chan<- T, done <-chan u.Never, value T) u.SelectResult {
	return u.Selection(u.SelectDone(done), u.SelectSend(valuesOut, value))
}

func defaultDrainErrorHandler(errorsOut chan<- error, done <-chan u.Never, err error) u.SelectResult {
	return u.Selection(u.SelectDone(done), u.SelectSend(errorsOut, err))
}
