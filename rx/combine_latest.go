package rx

import (
	"slices"
)

type combineLatestItem[T any] struct {
	subscription <-chan Message[T]
	unsubscribe  func()
	endOfStream  bool
	hasError     bool
	err          error
	hasValue     bool
	value        T
}

func CombineLatest[T any](sources ...Observable[T]) Observable[[]T] {
	return NewUnicastObservable(func(observer chan<- Message[[]T], done <-chan Void) {

		items := make([]*combineLatestItem[T], len(sources))

		for index, source := range sources {
			subscription, unsubscribe := source.Subscribe()
			defer unsubscribe()

			items[index] = &combineLatestItem[T]{subscription: subscription, unsubscribe: unsubscribe}
		}

		for {
			if combineLatestSelect(items, done) {
				return
			}

			if combineLatestHandleError(items, observer, done) {
				return
			}

			if combineLatestHandleResults(items, observer, done) {
				return
			}

			if combineLatestHandleEndOfStream(items) {
				return
			}
		}
	})
}

// combineLatestSelect returns isDone (ie: true for quit, false to keep going)
func combineLatestSelect[T any](items []*combineLatestItem[T], done <-chan Void) bool {

	var isDone bool

	selection := NewSelect()

	selection.Add(RevcItem(done, func(value Void, endOfStream bool) {
		isDone = true
	}))

	for item := range slices.Values(items) {

		item := item

		if !item.endOfStream {

			selection.Add(RevcItem(item.subscription, func(message Message[T], endOfStream bool) {

				if endOfStream {
					item.endOfStream = true
					item.unsubscribe()
					return
				}

				if message.IsError() {
					item.hasError = true
					item.err = message.Error
					return
				}

				if message.IsValue() {
					item.hasValue = true
					item.value = message.Value
					return
				}

				// Ignore unknown message types
			}))
		}
	}

	selection.Select()
	return isDone
}

// combineLatestHandleError returns isDone (ie: true for quit, false to keep going)
func combineLatestHandleError[T any](items []*combineLatestItem[T], observer chan<- Message[[]T], done <-chan Void) bool {

	errors := []error{}

	for item := range slices.Values(items) {
		if item.hasError {
			Append(&errors, item.err)
		}
	}

	if len(errors) == 0 {
		return false
	}

	for err := range slices.Values(errors) {
		if !Send1(NewError[[]T](err), observer, done) {
			return true
		}
	}

	return true
}

// combineLatestHandleResults returns isDone (ie: true for quit, false to keep going)
func combineLatestHandleResults[T any](items []*combineLatestItem[T], observer chan<- Message[[]T], done <-chan Void) bool {

	values := make([]T, len(items))

	for index, item := range items {
		if !item.hasValue {
			return false
		}
		values[index] = item.value
	}

	return !Send1(NewValue(values), observer, done)
}

// combineLatestHandleEndOfStream returns isDone (ie: true for quit, false to keep going)
func combineLatestHandleEndOfStream[T any](items []*combineLatestItem[T]) bool {

	for item := range slices.Values(items) {
		if !item.endOfStream {
			return false
		}
	}

	return true
}
