package rx

import (
	"slices"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type ReplaySubject[T any] struct {
	lock        *sync.Mutex
	subscribers *subscriberList[T]
	window      int
	history     []T
}

var _ Subject[any] = (*ReplaySubject[any])(nil)

func NewReplaySubject[T any](window int) *ReplaySubject[T] {

	var lock sync.Mutex

	return &ReplaySubject[T]{
		lock:        &lock,
		subscribers: NewSubscriberList[T](&lock),
		window:      window,
		history:     []T{},
	}
}

func (x *ReplaySubject[T]) Next(value T) bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	if !x.subscribers.Next(value) {
		return false
	}

	u.Append(&x.history, value)
	x.history = x.history[max(0, len(x.history)-x.window):len(x.history)]
	return true
}

func (x *ReplaySubject[T]) EndOfStream() {

	x.lock.Lock()
	defer x.lock.Unlock()

	x.subscribers.EndOfStream()
}

func (x *ReplaySubject[T]) Subscribe() (<-chan T, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	var unsubscribedCleanup, downstreamCleanup u.Event

	unsubscribed := u.NewChannel[u.Never](&unsubscribedCleanup, 0)
	downstream := u.NewChannel[T](&downstreamCleanup, 0)

	var initial []messageValue[T]

	for history := range slices.Values(x.history) {
		u.Append(&initial, messageValue[T]{value: history})
	}

	if x.subscribers.IsEndOfStream() {
		u.Append(&initial, messageValue[T]{endOfStream: true})
	}

	x.subscribers.AddSubscriber(downstream, unsubscribed, &downstreamCleanup, &unsubscribedCleanup, initial)

	return downstream, unsubscribedCleanup.Emit
}

func (x *ReplaySubject[T]) OnEndOfStream(listener func()) func() {
	x.lock.Lock()
	defer x.lock.Unlock()
	return x.subscribers.OnEndOfStream(listener)
}
