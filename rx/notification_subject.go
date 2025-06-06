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
}

var _ Observable[Void] = (*NotificationSubject)(nil)

func NewNotificationSubject(env *RxEnvironment) *NotificationSubject {

	var lock sync.Mutex

	return &NotificationSubject{
		env:         env,
		lock:        &lock,
		subscribers: NewSubscriberList[Void](env, &lock),
	}
}

func (x *NotificationSubject) Dispose() {

	x.lock.Lock()
	defer x.lock.Unlock()

	x.env.Error(x.subscribers.EndOfStream())
}

func (x *NotificationSubject) IsSignalled() bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	return x.subscribers.HasLatestValue()
}

func (x *NotificationSubject) Signal() error {

	x.lock.Lock()
	defer x.lock.Unlock()

	// Do nothjing if already signalled (notification subjects are idempotent)
	if x.subscribers.HasLatestValue() {
		return nil
	}

	if x.subscribers.IsEndOfStream() {
		return ErrEndOfStream
	}

	if err := x.subscribers.Next(Void{}); err != nil {
		return err
	}

	if err := x.subscribers.EndOfStream(); err != nil {
		return err
	}

	return nil
}

func (x *NotificationSubject) Subscribe(env *RxEnvironment) (<-chan Void, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	downstream, sendDownstreamEndOfStream := NewChannel[Void](env, 0)
	downstreamUnsubscribed, triggerDownstreamUnsubscribed := NewChannel[u.Never](env, 0)

	if latestValue, hasLatestValue := x.subscribers.GetLatestValue(); hasLatestValue {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed, latestValue)
	} else if x.subscribers.IsEndOfStream() {
		sendValuesThenEndOfStreamAsync(env, downstream, downstreamUnsubscribed)
	} else {
		x.subscribers.AddSubscriber(downstream, downstreamUnsubscribed, sendDownstreamEndOfStream, triggerDownstreamUnsubscribed, nil)
	}

	return downstream, triggerDownstreamUnsubscribed.Resolve
}
