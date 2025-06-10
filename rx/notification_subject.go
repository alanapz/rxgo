package rx

import (
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type Void = struct{}

type NotificationSubject struct {
	env         *RxEnvironment
	lock        *sync.Mutex
	subscribers *subscriberList[Void]
	disposed    bool
	signalled   bool
}

var _ Observable[Void] = (*NotificationSubject)(nil)
var _ Disposable = (*NotificationSubject)(nil)

func NewNotificationSubject(env *RxEnvironment) *NotificationSubject {

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

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.disposed
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

func (x *NotificationSubject) Subscribe(env *RxEnvironment) (<-chan Void, func(), error) {

	x.lock.Lock()
	defer x.lock.Unlock()

	downstream, sendDownstreamEndOfStream := NewChannel[Void](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	if x.signalled {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed, Void{})
	} else if x.disposed {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, nil)
	}

	return downstream, triggerDownstreamUnsubscribed.Emit
}
