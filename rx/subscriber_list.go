package rx

import (
	"maps"
	"slices"
	"sync"
	"sync/atomic"

	u "alanpinder.com/rxgo/v2/utils"
)

type subscriberId = uint64

type subscriberList[T any] struct {
	lock             *sync.Mutex
	nextSubscriberId atomic.Uint64
	valueSubscribers map[subscriberId]*messageQueue[T]
	errorSubscribers map[subscriberId]*messageQueue[error]
}

func NewSubscriberList[T any](lock *sync.Mutex) *subscriberList[T] {
	return &subscriberList[T]{
		lock:             lock,
		valueSubscribers: map[subscriberId]*messageQueue[T]{},
		errorSubscribers: map[subscriberId]*messageQueue[error]{},
	}
}

func (x *subscriberList[T]) PostValue(value T) {

	u.AssertLocked(x.lock)

	for valueSubscriber := range maps.Values(x.valueSubscribers) {
		valueSubscriber.Post(value)
	}
}

func (x *subscriberList[T]) ValuesComplete() {

	u.AssertLocked(x.lock)

	for valueSubscriber := range maps.Values(x.valueSubscribers) {
		valueSubscriber.PostComplete()
	}
}

func (x *subscriberList[T]) PostError(value error) {

	u.AssertLocked(x.lock)

	for errorSubscriber := range maps.Values(x.errorSubscribers) {
		errorSubscriber.Post(value)
	}
}

func (x *subscriberList[T]) ErrorsComplete() {

	u.AssertLocked(x.lock)

	for errorSubscriber := range maps.Values(x.errorSubscribers) {
		errorSubscriber.PostComplete()
	}
}

func (x *subscriberList[T]) AddSubscriber(values chan<- T, errors chan<- error, aborted <-chan u.Never, valuesCleanup *u.Cleanup, errorsCleanup *u.Cleanup, unsubscribedCleanup *u.Cleanup) {

	u.AssertLocked(x.lock)

	subscriberId := x.nextSubscriberId.Add(1)

	valuesCleanup.Add(func() {
		u.AssertLocked(x.lock)
		delete(x.valueSubscribers, subscriberId)
	})

	errorsCleanup.Add(func() {
		u.AssertLocked(x.lock)
		delete(x.errorSubscribers, subscriberId)
	})

	for cleanup := range slices.Values(u.Of(valuesCleanup, errorsCleanup)) {
		cleanup.AddPostCleanupHook(func() {

			u.AssertLocked(x.lock)

			if valuesCleanup.IsComplete() && errorsCleanup.IsComplete() {
				unsubscribedCleanup.Do()
			}
		})
	}

	x.valueSubscribers[subscriberId] = NewMessageQueue(x.lock, values, aborted, valuesCleanup)
	x.errorSubscribers[subscriberId] = NewMessageQueue(x.lock, errors, aborted, errorsCleanup)
}
