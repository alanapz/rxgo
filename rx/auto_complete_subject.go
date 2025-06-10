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
	disposed    bool
}

var _ Subject[any] = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any](env *RxEnvironment) *AutoCompleteSubject[T] {

	var lock sync.Mutex

	subject := &AutoCompleteSubject[T]{
		env:         env,
		lock:        &lock,
		subscribers: NewSubscriberList[T](env, &lock),
	}

	env.AddTryCleanup(subject.EndOfStream)
	return subject
}

func (x *AutoCompleteSubject[T]) IsDisposed() bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.disposed
}

func (x *AutoCompleteSubject[T]) Dispose() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return nil
	}

	x.disposed = true

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}

func (x *AutoCompleteSubject[T]) Next(values ...T) error {

	// AutoCompleteSubject only supports one value
	if len(values) > 1 {
		return fmt.Errorf("AutoCompleteSubject only supports a single value - additional values will be dropped")
	}

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return ErrDisposed
	}

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

	if x.disposed {
		return ErrDisposed
	}

	return x.subscribers.EndOfStream()
}

func (x *AutoCompleteSubject[T]) Subscribe(env *RxEnvironment) (<-chan T, func(), error) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return nil, nil, ErrDisposed
	}

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	endOfStream := x.subscribers.IsEndOfStream()
	latestValue, hasValue := x.subscribers.GetLatestValue()

	if endOfStream && hasValue {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed, latestValue)
	} else if endOfStream {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed)
	} else if err := x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed); err != nil {
		return nil, nil, err
	}

	return downstream, triggerDownstreamUnsubscribed.Emit, nil
}

func (x *AutoCompleteSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(x.env, source, x, endOfStreamPropagation)
}

func (x *AutoCompleteSubject[T]) OnEndOfStream(listener func()) func() {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.OnEndOfStream(listener)
}
