package rx

import (
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type AutoCompleteSubject[T any] struct {
	lock        *sync.Mutex
	subscribers *subscriberList[T]
}

var _ Subject[any] = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any]() *AutoCompleteSubject[T] {

	var lock sync.Mutex

	return &AutoCompleteSubject[T]{
		lock:        &lock,
		subscribers: NewSubscriberList[T](&lock),
	}
}

func (x *AutoCompleteSubject[T]) Next(value T) bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	if !x.subscribers.Next(value) {
		return false
	}

	x.subscribers.EndOfStream()

	return true
}

func (x *AutoCompleteSubject[T]) EndOfStream() {

	x.lock.Lock()
	defer x.lock.Unlock()

	x.subscribers.EndOfStream()
}

func (x *AutoCompleteSubject[T]) Subscribe() (<-chan T, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	var unsubscribedCleanup, downstreamCleanup u.Event

	unsubscribed := u.NewChannel[u.Never](&unsubscribedCleanup, 0)
	downstream := u.NewChannel[T](&downstreamCleanup, 0)

	var initial []messageValue[T]

	if latestValue, hasValue := x.subscribers.GetLatestValue(); hasValue {
		u.Append(&initial, messageValue[T]{value: latestValue})
	}

	if x.subscribers.IsEndOfStream() {
		u.Append(&initial, messageValue[T]{endOfStream: true})
	}

	x.subscribers.AddSubscriber(downstream, unsubscribed, &downstreamCleanup, &unsubscribedCleanup, initial)

	return downstream, unsubscribedCleanup.Emit
}

func (x *AutoCompleteSubject[T]) OnEndOfStream(listener func()) func() {
	x.lock.Lock()
	defer x.lock.Unlock()
	return x.subscribers.OnEndOfStream(listener)
}
