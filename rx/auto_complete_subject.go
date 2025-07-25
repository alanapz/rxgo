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
	ctx         *Context
	lock        *sync.Mutex
	subscribers *subscriberList[T]
	d           *ContextDisposable
}

var _ Subject[any] = (*AutoCompleteSubject[any])(nil)
var _ Disposable = (*AutoCompleteSubject[any])(nil)

func NewAutoCompleteSubject[T any](ctx *Context) (*AutoCompleteSubject[T], error) {

	var lock sync.Mutex

	var disposable *ContextDisposable

	if err := u.Wrap(NewContextDisposable(ctx, &lock))(&disposable); err != nil {
		return nil, err
	}

	subject := &AutoCompleteSubject[T]{
		ctx:         ctx,
		lock:        &lock,
		subscribers: NewSubscriberList[T](ctx, &lock),
		d:           disposable,
	}

	defer u.Lock(&lock)()

	if err := subject.d.AddCleanup(subject.cleanup); err != nil {
		return nil, err
	}

	return subject, nil
}

func (x *AutoCompleteSubject[T]) IsDisposed() bool {
	return x.d.IsDisposed()
}

func (x *AutoCompleteSubject[T]) IsDisposalInProgress() bool {
	return x.d.IsDisposalInProgress()
}

func (x *AutoCompleteSubject[T]) Dispose() {
	x.d.Dispose()
}

func (x *AutoCompleteSubject[T]) OnDisposalComplete() <-chan Void {
	return x.d.OnDisposalComplete()
}

func (x *AutoCompleteSubject[T]) Next(values ...T) error {

	// AutoCompleteSubject only supports one value
	if len(values) > 1 {
		return fmt.Errorf("AutoCompleteSubject only supports a single value - additional values will be dropped")
	}

	defer u.Lock(x.lock)()

	if x.IsDisposed() {
		return u.ErrDisposed
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

	defer u.Lock(x.lock)()

	if x.IsDisposed() {
		return u.ErrDisposed
	}

	return x.subscribers.EndOfStream()
}

func (x *AutoCompleteSubject[T]) Subscribe(ctx *Context) (<-chan T, func(), error) {

	defer u.Lock(x.lock)()

	if x.IsDisposed() {
		return nil, nil, u.ErrDisposed
	}

	var downstream chan T
	var sendDownstreamEndOfStream func()

	if err := u.Wrap2(NewChannel[T](ctx, 0))(&downstream, &sendDownstreamEndOfStream); err != nil {
		return nil, nil, err
	}

	var downstreamUnsubscribed chan u.Never
	var triggerDownstreamUnsubscribed func()

	if err := u.Wrap2(NewChannel[u.Never](ctx, 0))(&downstreamUnsubscribed, &triggerDownstreamUnsubscribed); err != nil {
		return nil, nil, err
	}

	endOfStream := x.subscribers.IsEndOfStream()
	latestValue, hasValue := x.subscribers.GetLatestValue()

	if endOfStream && hasValue {
		sendValuesThenEndOfStreamAsync(ctx, downstream, downstreamUnsubscribed, latestValue)
	} else if endOfStream {
		sendValuesThenEndOfStreamAsync(ctx, downstream, downstreamUnsubscribed)
	} else if err := x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed); err != nil {
		return nil, nil, err
	}

	return downstream, triggerDownstreamUnsubscribed, nil
}

func (x *AutoCompleteSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(x.ctx, source, x, endOfStreamPropagation)
}

func (x *AutoCompleteSubject[T]) OnEndOfStream(listener func()) func() {

	defer u.Lock(x.lock)()

	return x.subscribers.OnEndOfStream(listener)
}

func (x *AutoCompleteSubject[T]) cleanup() error {

	u.AssertLocked(x.lock)

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}
