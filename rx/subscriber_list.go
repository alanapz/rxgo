package rx

import (
	"maps"
	"sync"
	"sync/atomic"

	u "alanpinder.com/rxgo/v2/utils"
)

type subscriberId = uint64

type subscriberList[T any] struct {
	lock             *sync.Mutex
	nextSubscriberId atomic.Uint64
	subscribers      map[subscriberId]*messageQueue[T]
	endOfStream      bool
	onEndOfStream    u.Event
	hasValue         bool
	latestValue      T
}

func NewSubscriberList[T any](lock *sync.Mutex) *subscriberList[T] {
	return &subscriberList[T]{
		lock:        lock,
		subscribers: map[subscriberId]*messageQueue[T]{},
	}
}

func (x *subscriberList[T]) GetLatestValue() (T, bool) {
	u.AssertLocked(x.lock)
	return x.latestValue, x.hasValue
}

func (x *subscriberList[T]) IsEndOfStream() bool {
	u.AssertLocked(x.lock)
	return x.endOfStream
}

func (x *subscriberList[T]) Next(value T) bool {

	u.AssertLocked(x.lock)

	if x.endOfStream {
		return false
	}

	x.hasValue = true
	x.latestValue = value

	for subscriber := range maps.Values(x.subscribers) {
		subscriber.Next(value)
	}

	return true
}

func (x *subscriberList[T]) EndOfStream() {

	u.AssertLocked(x.lock)

	if x.endOfStream {
		return
	}

	x.endOfStream = true

	for subscriber := range maps.Values(x.subscribers) {
		subscriber.EndOfStream()
	}

	x.onEndOfStream.Emit()
}

func (x *subscriberList[T]) AddSubscriber(downstream chan<- T, aborted <-chan u.Never, downstreamCleanup *u.Event, unsubscribedCleanup *u.Event, initial []messageValue[T]) {

	u.AssertLocked(x.lock)

	subscriberId := x.nextSubscriberId.Add(1)

	downstreamCleanup.Add(func() {
		u.AssertLocked(x.lock)
		delete(x.subscribers, subscriberId)
	})

	downstreamCleanup.AddPostExecutionHook(func() {
		u.AssertLocked(x.lock)
		unsubscribedCleanup.Emit()
	})

	x.subscribers[subscriberId] = NewMessageQueue(x.lock, downstream, aborted, downstreamCleanup, initial)
}

func (x *subscriberList[T]) OnEndOfStream(listener func()) func() {
	u.AssertLocked(x.lock)
	return x.onEndOfStream.Add(listener)
}
