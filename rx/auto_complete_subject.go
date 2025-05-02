package rx

// import (
// 	"sync"
// )

// type AutoCompleteSubject[T any] struct {
// 	lock        sync.Mutex
// 	subscribers *subscriberList[T]
// 	signalled   bool
// 	value       T
// 	err         err
// }

// var _ Subject[any] = (*AutoCompleteSubject[any])(nil)

// func NewAutoCompleteSubject[T any]() *AutoCompleteSubject[T] {

// 	return &AutoCompleteSubject[T]{
// 		subscribers: map[SubscriberId]*Subscriber[T]{},
// 	}
// }

// func (x *AutoCompleteSubject[T]) PostValue(value T) {
// 	x.next(&value, nil)
// }

// func (x *AutoCompleteSubject[T]) PostError(err error) {
// 	x.next(nil, err)
// }

// func (x *AutoCompleteSubject[T]) PostComplete() {
// 	x.next(nil, nil)
// }

// func (x *AutoCompleteSubject[T]) next(value *T, err error) {

// 	x.lock.Lock()
// 	defer x.lock.Unlock()

// 	if x.signalled {
// 		panic("Already signalled")
// 	}

// 	x.signalled = true

// 	if value != nil {
// 		x.value = *value
// 		x.subscribers.PostValue(*value)
// 	}

// 	if err != nil {
// 		x.err = err
// 		x.subscribers.PostError(err)
// 	}

// 	x.subscribers.PostComplete()
// }

// func (x *AutoCompleteSubject[T]) Subscribe() (<-chan T, <-chan error, func()) {

// 	x.lock.Lock()
// 	defer x.lock.Unlock()

// 	values, cleanupValues := NewChannel[T](0)
// 	errors, cleanupErrors := NewChannel[error](0)
// 	done, cleanupDone := NewChannel[Never](0)

// 	if x.event != nil {

// 		// If we are already signalled, simply spawn a goroutine to send existing value to new channel
// 		// We can't simply send value as will deadlock
// 		go func() {

// 			defer cleanupValues()
// 			defer cleanupErrors()

// 			isValue, value := x.event.IsValue()
// 			isError, err := x.event.IsError()

// 			if isValue && Selection(SelectDone(done), SelectSend(values, value)) == DoneResult {
// 				return
// 			}

// 			if isError && Selection(SelectDone(done), SelectSend(errors, err)) == DoneResult {
// 				return
// 			}
// 		}()

// 		return values, errors, cleanupDone
// 	}

// 	subscriber := NewSubscriber(&x.lock, values, errors, done, func() {
// 		defer cleanupValues()
// 		defer cleanupErrors()
// 	})

// 	x.subscribers[x.nextSubscriberId.Add(1)] = subscriber

// 	return values, errors, cleanupDone
// }
