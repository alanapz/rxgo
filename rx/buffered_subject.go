package rx

import (
	"errors"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

func BufferUntil[T any](notifiers ...Observable[Void]) OperatorFunction[T, T] {
	return func(source Observable[T]) Observable[T] {
		return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

			var subject *BufferedSubject[T]

			if err := u.Wrap(NewBufferedSubject[T](ctx, notifiers))(&subject); err != nil {
				return err
			}

			defer subject.Dispose()

			subject.AddSource(source, PropogateEndOfStream)

			return drainObservable(drainObservableArgs[T]{
				Context:                ctx,
				Source:                 subject,
				Downstream:             downstream,
				DownstreamUnsubscribed: downstreamUnsubscribed,
			})
		})
	}
}

type BufferedSubject[T any] struct {
	ctx             *Context
	d               *ContextDisposable
	notifiers       []Observable[Void]
	lock            *sync.Mutex
	subscribers     *subscriberList[T]
	startWorkerOnce sync.Once
	buffer          []T
	signalled       bool // Whether we have been notified: if true,. pass messages to subscribers
	endOfStream     bool
}

var _ Subject[any] = (*BufferedSubject[any])(nil)
var _ Disposable = (*BufferedSubject[any])(nil)

type BufferedSubjectArgs struct {
	Context   *Context
	Notifiers []Observable[Void]
}

func NewBufferedSubject[T any](ctx *Context, notifiers []Observable[Void]) (*BufferedSubject[T], error) {

	var lock sync.Mutex

	var disposable *ContextDisposable

	if err := u.Wrap(NewContextDisposable(ctx, &lock))(&disposable); err != nil {
		return nil, err
	}

	subject := &BufferedSubject[T]{
		ctx:         ctx,
		d:           disposable,
		notifiers:   notifiers,
		lock:        &lock,
		subscribers: NewSubscriberList[T](ctx, &lock),
	}

	defer u.Lock(&lock)()

	if err := subject.d.AddCleanup(subject.cleanup); err != nil {
		return nil, err
	}

	subject.startWorker()
	return subject, nil
}

func (x *BufferedSubject[T]) IsDisposed() bool {
	return x.d.IsDisposed()
}

func (x *BufferedSubject[T]) IsDisposalInProgress() bool {
	return x.d.IsDisposalInProgress()
}

func (x *BufferedSubject[T]) Dispose() {
	x.d.Dispose()
}

func (x *BufferedSubject[T]) OnDisposalComplete() <-chan Void {
	return x.d.OnDisposalComplete()
}

func (x *BufferedSubject[T]) NotifySubscribers() error {

	defer u.Lock(x.lock)()

	if x.signalled {
		return nil
	}

	if x.IsDisposed() {
		return u.ErrDisposed
	}

	x.signalled = true
	return x.flushIfNecessary()
}

func (x *BufferedSubject[T]) Next(values ...T) error {

	defer u.Lock(x.lock)()

	if x.endOfStream {
		return ErrEndOfStream
	}

	if x.IsDisposed() {
		return u.ErrDisposed
	}

	u.Append(&x.buffer, values...)
	return x.flushIfNecessary()
}

func (x *BufferedSubject[T]) EndOfStream() error {

	defer u.Lock(x.lock)()

	if x.endOfStream {
		return nil
	}

	if x.IsDisposed() {
		return u.ErrDisposed
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

func (x *BufferedSubject[T]) Subscribe(ctx *Context) (<-chan T, func(), error) {

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

	if x.signalled && x.endOfStream {
		sendValuesThenEndOfStreamAsync(ctx, downstream, downstreamUnsubscribed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed)
	}

	return downstream, triggerDownstreamUnsubscribed, nil
}

func (x *BufferedSubject[T]) AddSource(source Observable[T], endOfStreamPropagation EndOfStreamPropagationPolicy) {
	PublishTo(x.ctx, source, x, endOfStreamPropagation)
}

func (x *BufferedSubject[T]) OnEndOfStream(listener func()) func() {

	defer u.Lock(x.lock)()

	return x.subscribers.OnEndOfStream(listener)
}

func (x *BufferedSubject[T]) startWorker() {

	x.startWorkerOnce.Do(func() {

		GoRun(x.ctx, func() {

			var upstream <-chan Void
			var unsubscribeFromUpstream func()

			err := u.Wrap2(Pipe(Merge(x.notifiers...), TakeUntil[Void](FromChannel(x.d.OnDisposalComplete)), First[Void]()).Subscribe(x.ctx))(&upstream, &unsubscribeFromUpstream)

			// XXXX
			// if errors.Is(err, ErrDone) {
			// 	return
			// }

			if err != nil {
				x.ctx.Error(err)
				return
			}

			defer unsubscribeFromUpstream()

			for range upstream {
				x.ctx.Assert(x.NotifySubscribers())
			}
		})
	})
}

func (x *BufferedSubject[T]) cleanup() error {

	u.AssertLocked(x.lock)

	if !x.signalled {
		return errors.New("dispose before subscribers notified: buffered events will be lost")
	}

	if !x.endOfStream {
		return errors.New("dispose before end-of-stream: pending events will be lost")
	}

	return nil
}
