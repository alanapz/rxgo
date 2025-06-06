package rx

import (
	"errors"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func BufferUntil[T any](notifiers ...Observable[Void]) OperatorFunction[T, T] {

	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(args UnicastObserverArgs[T]) {

			subject := NewBufferedSubject[T](BufferedSubjectArgs{Environment: args.Environment, Notifiers: notifiers})
			defer subject.Dispose()

			subject.AddSource(source, PropogateEndOfStream)

			drainObservable(drainObservableArgs[T]{
				Environment:            args.Environment,
				Source:                 subject,
				Downstream:             args.Downstream,
				DownstreamUnsubscribed: args.DownstreamUnsubscribed,
			})
		})
	}
}

type BufferedSubject[T any] struct {
	env               *RxEnvironment
	notifiers         []Observable[Void]
	lock              *sync.Mutex
	subscribers       *subscriberList[T]
	notifySubscribers bool // Whether we have been notified: if true,. pass messages to subscribers
	buffer            []T
	endOfStream       bool
	disposed          bool
	dispose           *NotificationSubject
	startWorkerOnce   sync.Once
}

var _ Subject[any] = (*BufferedSubject[any])(nil)

type BufferedSubjectArgs struct {
	Environment *RxEnvironment
	Notifiers   []Observable[Void]
}

func NewBufferedSubject[T any](args BufferedSubjectArgs) *BufferedSubject[T] {

	var lock sync.Mutex

	subject := &BufferedSubject[T]{
		env:         args.Environment,
		notifiers:   args.Notifiers,
		lock:        &lock,
		subscribers: NewSubscriberList[T](args.Environment, &lock),
		dispose:     NewNotificationSubject(args.Environment),
	}

	subject.startWorker()

	return subject
}

func (x *BufferedSubject[T]) Dispose() {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return
	}

	if !x.notifySubscribers {
		x.env.Error(errors.New("dispose before subscribers notified: buffered events will be lost"))
		return
	}

	if !x.endOfStream {
		x.env.Error(errors.New("dispose before end-of-stream: pending events will be lost"))
		return
	}

	x.disposed = true

	x.env.Error(x.dispose.Signal())
}

func (x *BufferedSubject[T]) NotifySubscribers() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.notifySubscribers {
		return nil
	}

	if x.disposed {
		return ErrDisposed
	}

	x.notifySubscribers = true
	return x.flushIfNecessary()
}

func (x *BufferedSubject[T]) Next(values ...T) error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.endOfStream {
		return ErrEndOfStream
	}

	if x.disposed {
		return ErrDisposed
	}

	u.Append(&x.buffer, values...)
	return x.flushIfNecessary()
}

func (x *BufferedSubject[T]) EndOfStream() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.endOfStream {
		return nil
	}

	if x.disposed {
		return ErrDisposed
	}

	x.endOfStream = true
	return x.flushIfNecessary()
}

func (x *BufferedSubject[T]) flushIfNecessary() error {

	u.AssertLocked(x.lock)

	if !x.notifySubscribers {
		return nil
	}

	if err := x.subscribers.Next(x.buffer...); err != nil { // Safe if buffer is empty
		return err
	}

	x.buffer = nil

	if x.endOfStream {
		if err := x.subscribers.EndOfStream(); err != nil { //Safe to be called multiple times
			return err
		}
	}

	return nil
}

func (x *BufferedSubject[T]) Subscribe(env *RxEnvironment) (<-chan T, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		x.env.Error(ErrDisposed)
		return nil, nil
	}

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	if x.subscribers.IsEndOfStream() {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, nil)
	}

	return downstream, triggerDownstreamUnsubscribed.Resolve
}

func (x *BufferedSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {

	x.lock.Lock()
	defer x.lock.Unlock()

	PublishTo(PublishToArgs[T]{Environment: x.env, Source: source, Sink: x, PropogateEndOfStream: bool(endOfStreamPropagation)})
}

func (x *BufferedSubject[T]) OnEndOfStream(listener func()) func() {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.OnEndOfStream(listener)
}

func (x *BufferedSubject[T]) startWorker() {

	x.startWorkerOnce.Do(func() {

		x.env.Execute(func() {

			upstream, unsubscribeFromUpstream := Pipe(Merge(x.notifiers...), TakeUntil[Void](x.dispose), First[Void]()).Subscribe(x.env)
			defer unsubscribeFromUpstream()

			for range upstream {
				x.env.Error(x.NotifySubscribers())
			}
		})
	})
}
