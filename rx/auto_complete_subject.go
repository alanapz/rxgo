package rx

import (
	"maps"
	"sync"
	"sync/atomic"
)

type AutoCompleteSubject[T any] struct {
	lock             sync.Mutex
	nextSubscriberId atomic.Uint64
	subscribers      map[SubscriberId]*Subscriber[T]
	event            *SubscriberEvent[T]
}

var _ Subject[any] = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any]() *AutoCompleteSubject[T] {

	return &AutoCompleteSubject[T]{
		subscribers: map[SubscriberId]*Subscriber[T]{},
	}
}

func (x *AutoCompleteSubject[T]) Value(value T) {
	x.next(NewValueEvent(value))
}

func (x *AutoCompleteSubject[T]) Error(err error) {
	x.next(NewErrorEvent[T](err))
}

func (x *AutoCompleteSubject[T]) Complete() {
	x.next(NewCompleteEvent[T]())
}

func (x *AutoCompleteSubject[T]) next(event SubscriberEvent[T]) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.event != nil {
		panic("Already signalled")
	}

	x.event = &event

	x.broadcastEvent(event)

	if !event.IsComplete() {
		x.broadcastEvent(NewCompleteEvent[T]())
	}
}
func (x *AutoCompleteSubject[T]) broadcastEvent(event SubscriberEvent[T]) {

	AssertLocked(&x.lock)

	for subscriber := range maps.Values(x.subscribers) {
		subscriber.PostEvent(event)
	}
}

func (x *AutoCompleteSubject[T]) Subscribe() (<-chan T, <-chan error, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	values, cleanupValues := NewChannel[T](0)
	errors, cleanupErrors := NewChannel[error](0)
	done, cleanupDone := NewChannel[Never](0)

	if x.event != nil {

		// If we are already signalled, simply spawn a goroutine to send existing value to new channel
		// We can't simply send value as will deadlock
		go func() {

			defer cleanupValues()
			defer cleanupErrors()

			isValue, value := x.event.IsValue()
			isError, err := x.event.IsError()

			if isValue && Selection(SelectDone(done), SelectSend(values, value)) {
				return
			}

			if isError && Selection(SelectDone(done), SelectSend(errors, err)) {
				return
			}
		}()

		return values, errors, cleanupDone
	}

	subscriber := NewSubscriber(&x.lock, values, errors, done, func() {
		defer cleanupValues()
		defer cleanupErrors()
	})

	x.subscribers[x.nextSubscriberId.Add(1)] = subscriber

	return values, errors, cleanupDone
}

func (x *AutoCompleteSubject[T]) Pipe(operator OperatorFunction[T, T]) Observable[T] {
	return operator(x)
}
