package rx

// import (
// 	"slices"
// )

// type combineLatestCtx[T any] struct {
// 	ValuesOut chan<- []T
// 	ErrorsOut chan<- error
// 	Done      <-chan Never
// 	Items     []*combineLatestItem[T]
// }

// type combineLatestItem[T any] struct {
// 	ValuesIn    <-chan T
// 	ErrorsIn    <-chan error
// 	Unsubscribe func()
// 	Disposed    bool
// 	Emitted     bool
// 	LastEmitted T
// 	//
// 	IncomingValue SelectReceiveResult[T]
// 	IncomingError SelectReceiveResult[error]
// }

// func CombineLatest[T any](sources ...Observable[T]) Observable[[]T] {
// 	return NewUnicastObservable(func(valuesOut chan<- []T, errorsOut chan<- error, done <-chan Never) {

// 		ctx := combineLatestCtx[T]{
// 			ValuesOut: valuesOut,
// 			ErrorsOut: errorsOut,
// 			Done:      done,
// 			Items:     make([]*combineLatestItem[T], len(sources)),
// 		}

// 		for index, source := range sources {
// 			valuesIn, errorsIn, unsubscribe := source.Subscribe()
// 			defer unsubscribe()

// 			ctx.Items[index] = &combineLatestItem[T]{
// 				ValuesIn:    valuesIn,
// 				ErrorsIn:    errorsIn,
// 				Unsubscribe: unsubscribe,
// 			}
// 		}

// 		for {
// 			if ctx.SelectNext() {
// 				return
// 			}

// 			if ctx.HandleErrors() {
// 				return
// 			}

// 			if ctx.HandleResults() {
// 				return
// 			}

// 			if ctx.HandleEndOfStream() {
// 				return
// 			}
// 		}
// 	})
// }

// // combineLatestSelect returns isDone (ie: true for quit, false to keep going)
// func (x *combineLatestCtx[T]) SelectNext() bool {

// 	selections := []SelectItem{}

// 	Append(&selections, SelectDone(x.Done))

// 	for item := range slices.Values(x.Items) {
// 		if !item.Disposed {
// 			Append(&selections, SelectReceiveInto(item.ValuesIn, &item.IncomingValue))
// 			Append(&selections, SelectReceiveInto(item.ErrorsIn, &item.IncomingError))
// 		}
// 	}

// 	return Selection(selections...)
// }

// func (x *combineLatestCtx[T]) HandleErrors() bool {

// 	errors := []error{}

// 	for item := range slices.Values(x.Items) {
// 		if !item.Disposed && item.IncomingError.Valid {
// 			Append(&errors, item.IncomingError.Value)
// 		}
// 	}

// 	if len(errors) == 0 {
// 		return false
// 	}

// 	for err := range slices.Values(errors) {
// 		if Selection(SelectDone(x.Done), SelectSend(x.ErrorsOut, err)) {
// 			break
// 		}
// 	}

// 	return true
// }

// func (x *combineLatestCtx[T]) HandleResults() bool {

// 	for item := range slices.Values(x.Items) {
// 		if !item.Disposed && item.IncomingValue.Valid {
// 			item.Emitted = true
// 			item.LastEmitted = item.IncomingValue.Value
// 		}
// 	}

// 	values := make([]T, len(x.Items))

// 	for index, item := range x.Items {
// 		if !item.Emitted {
// 			return false
// 		}
// 		values[index] = item.LastEmitted
// 	}

// 	return Selection(SelectDone(x.Done), SelectSend(x.ValuesOut, values))
// }

// func (x *combineLatestCtx[T]) HandleEndOfStream() bool {

// 	for item := range slices.Values(x.Items) {
// 		if !item.Disposed && (item.IncomingValue.EndOfStream || item.IncomingError.EndOfStream) {
// 			item.Unsubscribe()
// 			item.ValuesIn = nil
// 			item.ErrorsIn = nil
// 			item.Disposed = true
// 		}
// 	}

// 	for item := range slices.Values(x.Items) {
// 		if !item.Disposed {
// 			return false
// 		}
// 	}

// 	return true
// }
