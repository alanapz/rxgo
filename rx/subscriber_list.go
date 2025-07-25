package rx

import (
	"errors"
	"maps"
	"slices"
	"sync"
	"sync/atomic"

	u "alanpinder.com/rxgo/v2/utils"
)

type subscriberId = uint64

type subscriberList[T any] struct {
	ctx              *Context
	lock             *sync.Mutex
	nextSubscriberId atomic.Uint64
	subscribers      map[subscriberId]*messageQueue[T]
	endOfStream      bool
	onEndOfStream    u.Event
	hasValue         bool
	latestValue      T
}

func NewSubscriberList[T any](ctx *Context, lock *sync.Mutex) *subscriberList[T] {
	return &subscriberList[T]{
		ctx:         ctx,
		lock:        lock,
		subscribers: map[subscriberId]*messageQueue[T]{},
	}
}

func (x *subscriberList[T]) HasLatestValue() bool {
	u.AssertLocked(x.lock)
	return x.hasValue
}

func (x *subscriberList[T]) GetLatestValue() (T, bool) {
	u.AssertLocked(x.lock)
	return x.latestValue, x.hasValue
}

func (x *subscriberList[T]) IsEndOfStream() bool {
	u.AssertLocked(x.lock)
	return x.endOfStream
}

func (x *subscriberList[T]) Next(values ...T) error {

	u.AssertLocked(x.lock)

	if len(values) == 0 {
		return nil
	}

	if x.endOfStream {
		return ErrEndOfStream
	}

	x.hasValue = true

	var postErrors []error

	for value := range slices.Values(values) {

		x.latestValue = value

		for subscriber := range maps.Values(x.subscribers) {
			if err := subscriber.Next(value); err != nil {
				u.Append(&postErrors, err)
			}
		}
	}

	return errors.Join(postErrors...)
}

func (x *subscriberList[T]) EndOfStream() error {

	u.AssertLocked(x.lock)

	if x.endOfStream {
		return nil
	}

	x.endOfStream = true

	var postErrors []error

	for subscriber := range maps.Values(x.subscribers) {
		if err := subscriber.EndOfStream(); err != nil {
			u.Append(&postErrors, err)
		}
	}

	x.onEndOfStream.Emit()
	return errors.Join(postErrors...)
}

func (x *subscriberList[T]) AddSubscriber(downstream chan<- T, aborted <-chan u.Never, downstreamCleanup *u.Event, unsubscribedCleanup *u.Event, initial ...messageValue[T]) error {

	u.AssertLocked(x.lock)

	if x.endOfStream {
		return ErrEndOfStream
	}

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
	return nil
}

func (x *subscriberList[T]) OnEndOfStream(listener func()) func() {
	u.AssertLocked(x.lock)

	if x.endOfStream {
		listener()
		return u.DoNothing
	}

	return x.onEndOfStream.Add(listener)
}
