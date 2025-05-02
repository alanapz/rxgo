package rx

// import (
// 	"errors"
// 	"fmt"
// 	"maps"
// 	"sync"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// type SubscriberId = uint64

// type subscriberList[T any] struct {
// 	lock             *sync.Mutex
// 	nextSubscriberId uint64
// 	subscribers      map[SubscriberId]*subscriber[T]
// }

// func NeSubscriberList[T any](lock *sync.Mutex) *subscriberList[T] {

// 	u.Assert(lock != nil)

// 	return &SubscriberList[T]{
// 		lock:        lock,
// 		subscribers: map[SubscriberId]*Subscriber[T]{},
// 	}
// }

// func (x *subscriberList[T]) PostValue(value T) {

// 	u.AssertLocked(x.lock)

// 	postErrors := map[SubscriberId]error{}

// 	for subscriberId, subscriber := range x.subscribers {
// 		if err := subscriber.PostValue(value); err != nil {
// 			postErrors[subscriberId] = err
// 		}
// 	}

// 	for subscriberId, postError := range postErrors {
// 		x.disposeIfNecessary(subscriberId, postError, x.valueSubscribers)
// 	}
// }

// func (x *subscriberList[T]) PostError(value error) {

// 	u.AssertLocked(x.lock)

// 	postError := map[SubscriberId]error{}

// 	for subscriberId, subscriber := range x.errorSubscribers {
// 		if err := subscriber.PostValue(value); err != nil {
// 			postError[subscriberId] = err
// 		}
// 	}

// 	for subscriberId, postError := range postError {
// 		x.disposeIfNecessary(subscriberId, postError, x.errorSubscribers)
// 	}
// }

// func (x *subscriberList[T]) PostComplete(value error) {

// 	u.AssertLocked(x.lock)

// 	postError := map[SubscriberId]error{}

// 	for subscriberId, subscriber := range x.errorSubscribers {
// 		if err := subscriber.PostValue(value); err != nil {
// 			postError[subscriberId] = err
// 		}
// 	}

// 	for subscriberId, postError := range postError {
// 		x.disposeIfNecessary(subscriberId, postError, x.errorSubscribers)
// 	}
// }

// func (x *SubscriberList[T]) disposeIfNecessary(subscriberId SubscriberId, postError error, target map[SubscriberId]*MessageQueue[T]) {

// 	isEndOfStream := errors.Is(err, ErrEndOfStream)
// 	isDone := errors.Is(err, ErrDisposed)

// 	if isEndOfStream || isDone {

// 		if existing, ok := x.valueSubscribers[subscriberId]; ok {
// 			delete(x.valueSubscribers, subscriberId) // Safe, if we arrived here means queue will clean itself up
// 			return
// 		}
// 	}

// 	if isDone {

// 	}

// 	println(fmt.Sprintf("removing %d, error: %v", s))

// 	if err == ErrEndOfStream {
// 		x
// 	}

// 	u.AssertLocked(x.lock)

// 	errors := map[*MessageQueue[T]]error{}

// 	for subscriber := range maps.Values(x.valueSubscribers) {
// 		if err := subscriber.PostValue(value); err != nil {
// 			errors[subscriber] = err
// 		}
// 	}
// }

// func (x *SubscriberList[T]) PostError(err error) {

// 	u.AssertLocked(x.lock)

// 	errors := map[*Subscriber[T]]error{}

// 	for subscriber := range maps.Values(x.subscribers) {
// 		if err := subscriber.PostValue(value); err != nil {
// 			errors[subscriber] = err
// 		}
// 	}
// }

// func (x *SubscriberList[T]) PostComplete() {

// 	u.AssertLocked(x.lock)

// 	errors := map[*Subscriber[T]]error{}

// 	for subscriber := range maps.Values(x.subscribers) {
// 		if err := subscriber.PostValue(value); err != nil {
// 			errors[subscriber] = err
// 		}
// 	}
// }

// func (x *SubscriberList[T]) AddSubscriber(lock *sync.Mutex, valuesOut chan<- T, errorsOut chan<- error, done <-chan u.Never, cleanupValues func(), cleanupErrors func()) {

// 	u.AssertLocked(x.lock)

// 	x.nextSubscriberId++
// 	x.subscribers[x.nextSubscriberId] = subscriber
// }
