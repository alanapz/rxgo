package rx

import (
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type Void = struct{}

type NotificationSubject struct {
	ctx         *Context
	lock        *sync.Mutex
	subscribers *subscriberList[Void]
	d           *ContextDisposable
	signalled   bool
}

var _ Observable[Void] = (*NotificationSubject)(nil)
var _ Disposable = (*NotificationSubject)(nil)

func NewNotificationSubject(ctx *Context) *NotificationSubject {

	var lock sync.Mutex

	subject := &NotificationSubject{
		env:         env,
		lock:        &lock,
		subscribers: NewSubscriberList[Void](env, &lock),
	}

	env.AddTryCleanup(subject.Dispose)
	return subject
}

func (x *NotificationSubject) IsDisposed() bool {
	return x.d.IsDisposed()
}
func (x *NotificationSubject) IsDisposalInProgress() bool {
	return x.d.IsDisposed()
}

func (x *NotificationSubject) Dispose() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	x.disposed = true

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}

func (x *NotificationSubject) IsSignalled() bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.signalled
}

func (x *NotificationSubject) Signal() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	// Do nothing if already signalled (notification subjects are idempotent)
	if x.signalled {
		return nil
	}

	if x.disposed {
		return ErrDisposed
	}

	x.signalled = true

	if err := x.subscribers.Next(Void{}); err != nil {
		return err
	}

	x.disposed = true

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}

func (x *NotificationSubject) Subscribe(ctx *Context) (<-chan Void, func(), error) {

	x.lock.Lock()
	defer x.lock.Unlock()

	var downstream chan Void
	var sendDownstreamEndOfStream func()

	if err := u.W2(NewChannel[Void](ctx, 0))(&downstream, &sendDownstreamEndOfStream); err != nil {
		return nil, nil, err
	}

	var downstreamUnsubscribed chan u.Never
	var triggerDownstreamUnsubscribed func()

	if err := u.W2(NewChannel[u.Never](ctx, 0))(&downstreamUnsubscribed, &triggerDownstreamUnsubscribed); err != nil {
		return nil, nil, err
	}

	if x.signalled {
		sendValuesThenEndOfStreamAsync(ctx, downstream, downstreamUnsubscribed, Void{})
	} else if x.disposed {
		sendValuesThenEndOfStreamAsync(ctx, downstream, downstreamUnsubscribed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, nil)
	}

	return downstream, triggerDownstreamUnsubscribed, nil
}

func (x *NotificationSubject) cleanup() error {

	u.AssertLocked(x.lock)

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}
