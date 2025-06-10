package rx

import (
	"slices"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type ReplaySubject[T any] struct {
	env         *RxEnvironment
	lock        *sync.Mutex
	subscribers *subscriberList[T]
	window      int
	history     []T
}

var _ Subject[any] = (*ReplaySubject[any])(nil)

func NewReplaySubject[T any](env *RxEnvironment, window int) *ReplaySubject[T] {

	var lock sync.Mutex

	return &ReplaySubject[T]{
		env:         env,
		lock:        &lock,
		subscribers: NewSubscriberList[T](env, &lock),
		window:      window,
	}
}

func (x *ReplaySubject[T]) Next(values ...T) error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if err := x.subscribers.Next(values...); err != nil {
		return err
	}

	u.Append(&x.history, values...)
	x.history = x.history[max(0, len(x.history)-x.window):len(x.history)]
	return nil
}

func (x *ReplaySubject[T]) EndOfStream() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.EndOfStream()
}

func (x *ReplaySubject[T]) Subscribe(env *RxEnvironment) (<-chan T, func(), error) {

	x.lock.Lock()
	defer x.lock.Unlock()

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	var initial []messageValue[T]

	for history := range slices.Values(x.history) {
		u.Append(&initial, messageValue[T]{value: history})
	}

	if x.subscribers.IsEndOfStream() {
		u.Append(&initial, messageValue[T]{endOfStream: true})
	}

	x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, initial)

	return downstream, triggerDownstreamUnsubscribed.Emit
}

func (x *ReplaySubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(x.env, source, x, endOfStreamPropagation)
}

func (x *ReplaySubject[T]) OnEndOfStream(listener func()) func() {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.OnEndOfStream(listener)
}
