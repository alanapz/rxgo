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

			subject.AddSource(source, PropogateEndOfStreamFromSourceToSink)

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
	env             *RxEnvironment
	notifiers       []Observable[Void]
	lock            *sync.Mutex
	subscribers     *subscriberList[T]
	buffer          []T
	signalled       bool // Whether we have been notified: if true,. pass messages to subscribers
	endOfStream     bool
	disposed        bool
	dispose         *NotificationSubject
	startWorkerOnce sync.Once
}

var _ Subject[any] = (*BufferedSubject[any])(nil)
var _ Disposable = (*BufferedSubject[any])(nil)

type BufferedSubjectArgs struct {
	Environment *RxEnvironment
	Notifiers   []Observable[Void]
}

func NewBufferedSubject[T any](args BufferedSubjectArgs) *BufferedSubject[T] {

	env, notifiers := args.Environment, args.Notifiers

	u.Require(env, notifiers)

	var lock sync.Mutex

	subject := &BufferedSubject[T]{
		env:         env,
		notifiers:   args.Notifiers,
		lock:        &lock,
		subscribers: NewSubscriberList[T](env, &lock),
		dispose:     NewNotificationSubject(env),
	}

	subject.startWorker()

	env.AddTryCleanup(subject.Dispose)
	return subject
}

func (x *BufferedSubject[T]) IsDisposed() bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.disposed
}

func (x *BufferedSubject[T]) Dispose() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return nil
	}

	if !x.signalled {
		return errors.New("dispose before subscribers notified: buffered events will be lost")
	}

	if !x.endOfStream {
		return errors.New("dispose before end-of-stream: pending events will be lost")
	}

	x.disposed = true

	if err := x.dispose.Signal(); err != nil {
		return err
	}

	return nil
}

func (x *BufferedSubject[T]) NotifySubscribers() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.signalled {
		return nil
	}

	if x.disposed {
		return ErrDisposed
	}

	x.signalled = true
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

	if !x.signalled {
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

func (x *BufferedSubject[T]) Subscribe(env *RxEnvironment) (<-chan T, func(), error) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return nil, nil, ErrDisposed
	}

	downstream, sendDownstreamEndOfStream := NewChannel[T](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	endOfStream := x.subscribers.IsEndOfStream()

	if x.signalled && x.endOfStream {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed)
	} else if x.disposed {
		x.env.Error(ErrDisposed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed)
	}

	return downstream, triggerDownstreamUnsubscribed.Emit
}

func (x *BufferedSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(x.env, source, x, endOfStreamPropagation)
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
