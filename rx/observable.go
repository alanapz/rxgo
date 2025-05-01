package rx

import (
	"alanpinder.com/rxgo/v2/ux"
)

type Observable[T any] interface {
	Subscribe() (<-chan T, <-chan error, func())
	Pipe(func(Observable[T]) Observable[T]) Observable[T]
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
	done                  <-chan Never
	additionalSelections  func() []SelectItem
	afterSelectionHandler func(*SelectReceiveMessage[T], *SelectReceiveMessage[error]) AfterSelectionResult
	valueHandler          func(chan<- T, <-chan Never, T) SelectResult
	errorHandler          func(chan<- error, <-chan Never, error) SelectResult
}

func drainObservable[T any](args drainObservableArgs[T]) SelectResult {

	source, valuesOut, errorsOut, done, additionalSelections, afterSelection, valueHandler, errorHandler :=
		args.source,
		args.valuesOut,
		args.errorsOut,
		args.done,
		MustCoalesce(args.additionalSelections, defaultAdditionalSelections),
		MustCoalesce(args.afterSelectionHandler, defaultAfterSelectionHandler),
		MustCoalesce(args.valueHandler, defaultDrainValueHandler),
		MustCoalesce(args.errorHandler, defaultDrainErrorHandler)

	valuesIn, errorsIn, unsubscribe := source.Subscribe()
	defer unsubscribe()

	for {

		if valuesIn == nil && errorsIn == nil {
			return ContinueResult
		}

		var valueMsg SelectReceiveMessage[T]
		var errorMsg SelectReceiveMessage[error]

		selectionItems := ux.Of(
			SelectDone(done),
			SelectReceive(&valuesIn, &valueMsg),
			SelectReceive(&errorsIn, &errorMsg),
		)

		Append(&selectionItems, additionalSelections()...)

		if Selection(selectionItems...) {
			return DoneResult
		}

		switch afterSelection(&valueMsg, &errorMsg) {
		case DropMessage:
			continue
		case ReturnContinue:
			return ContinueResult
		case ReturnDone:
			return DoneResult
		}

		if valueMsg.HasValue && valueHandler(valuesOut, done, valueMsg.Value) == DoneResult {
			return DoneResult
		}

		if errorMsg.HasValue && errorHandler(errorsOut, done, errorMsg.Value) == DoneResult {
			return DoneResult
		}
	}
}

func defaultAdditionalSelections() []SelectItem {
	return Zero[[]SelectItem]()
}

func defaultAfterSelectionHandler[T any](valueMsg *SelectReceiveMessage[T], errorMsg *SelectReceiveMessage[error]) AfterSelectionResult {
	return ContinueMessage
}

func defaultDrainValueHandler[T any](valuesOut chan<- T, done <-chan Never, value T) SelectResult {
	return Selection(SelectDone(done), SelectSend(valuesOut, value))
}

func defaultDrainErrorHandler(errorsOut chan<- error, done <-chan Never, err error) SelectResult {
	return Selection(SelectDone(done), SelectSend(errorsOut, err))
}
