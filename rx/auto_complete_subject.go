package rx

import (
	"fmt"
	"maps"
	"runtime/debug"
	"sync"
)

type SubscriptionId = int

type AutoCompleteSubject[T any] struct {
	lock               sync.Mutex
	nextSubscriptionId SubscriptionId
	subscriptions      map[SubscriptionId]*AutoCompleteSubscription[T]
	msg                *Message[T]
}

type AutoCompleteSubscription[T any] struct {
	SubscriptionId SubscriptionId
	Subscriber     *Channel[Message[T]]
	Unsubcribed    *Channel[Void]
}

var _ Observable[any] = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any]() *AutoCompleteSubject[T] {
	return &AutoCompleteSubject[T]{
		subscriptions: map[SubscriptionId]*AutoCompleteSubscription[T]{},
	}
}

func (x *AutoCompleteSubject[T]) Signal(value T) {
	x.doSignal(NewValue(value))
}

func (x *AutoCompleteSubject[T]) SignalError(err error) {
	x.doSignal(NewError[T](err))
}

func (x *AutoCompleteSubject[T]) doSignal(msg Message[T]) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.msg != nil {
		panic("Already signalled")
	}

	x.msg = &msg

	for {

		if len(x.subscriptions) == 0 {
			// Nothing to do
			return
		}

		selectBuilder := NewSelect()

		for subscription := range maps.Values(x.subscriptions) {

			subscription := subscription

			selectBuilder.Add(SendItem(subscription.Subscriber.Channel, *x.msg, func() {
				x.disposeSubscription(subscription)
			}))

			selectBuilder.Add(RevcItem(subscription.Unsubcribed.Channel, func(value Void, endOfStream bool) {
				x.disposeSubscription(subscription)
			}))
		}

		selectBuilder.Select()
	}
}

func (x *AutoCompleteSubject[T]) disposeSubscription(subscription *AutoCompleteSubscription[T]) {

	// lock must be held !

	delete(x.subscriptions, subscription.SubscriptionId)
	subscription.Subscriber.Dispose()
	subscription.Unsubcribed.Dispose()

}

func (x *AutoCompleteSubject[T]) Subscribe() (<-chan Message[T], func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	subscriber := NewChannel[Message[T]](fmt.Sprintf("%s: for: %s", debug.Stack(), "x.stack"))
	unsubscribed := NewChannel[Void](fmt.Sprintf("%s: for: %s", debug.Stack(), "x.stack"))

	if x.msg != nil {

		// If we are already signalled, simply spawn a goroutine to send existing value to new channel
		// We can't simply send value as will deadlock
		go func() {
			defer subscriber.Dispose()
			defer unsubscribed.Dispose()
			Send1(*x.msg, subscriber.Channel, unsubscribed.Channel)
		}()

		return subscriber.Channel, unsubscribed.Dispose
	}

	x.nextSubscriptionId++

	subscription := &AutoCompleteSubscription[T]{
		SubscriptionId: x.nextSubscriptionId,
		Subscriber:     subscriber,
		Unsubcribed:    unsubscribed,
	}

	x.subscriptions[subscription.SubscriptionId] = subscription

	return subscriber.Channel, unsubscribed.Dispose
}
