package rx

import (
	"fmt"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

/*
An AutoComplete subject onlt accept one value, and completes after single value.
*/

type AutoCompleteSubject[T any] struct {
	env         *RxEnvironment
	lock        *sync.Mutex
	subscribers *subscriberList[T]
}

var _ Subject[any] = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any](env *RxEnvironment) *AutoCompleteSubject[T] {

	var lock sync.Mutex

	return &AutoCompleteSubject[T]{
		env:         env,
		lock:        &lock,
		subscribers: NewSubscriberList[T](env, &lock),
	}
}

func (x *AutoCompleteSubject[T]) Next(values ...T) error {

	// AutoCompleteSubject only supports one value
	if len(values) > 1 {
		return fmt.Errorf("AutoCompleteSubject only supports a single value - additional values will be dropped")
	}

	x.lock.Lock()
	defer x.lock.Unlock()

	if err := x.subscribers.Next(values...); err != nil {
		return err
	}

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}

func (x *AutoCompleteSubject[T]) EndOfStream() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.EndOfStream()
}

func (x *AutoCompleteSubject[T]) Subscribe(env *RxEnvironment) (<-chan T, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	var initial []messageValue[T]

	if latestValue, hasValue := x.subscribers.GetLatestValue(); hasValue {
		u.Append(&initial, messageValue[T]{value: latestValue})
	}

	if x.subscribers.IsEndOfStream() {
		u.Append(&initial, messageValue[T]{endOfStream: true})
	}

	x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, initial)

	return downstream, triggerDownstreamUnsubscribed.Resolve
}

func (x *AutoCompleteSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(PublishToArgs[T]{Source: source, Sink: x, PropogateEndOfStream: bool(endOfStreamPropagation)})
}

func (x *AutoCompleteSubject[T]) OnEndOfStream(listener func()) func() {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.OnEndOfStream(listener)
}
