package main

// import "sync"

// type Latch[T any] struct {
// 	Lock        sync.Mutex
// 	Signalled   bool
// 	Value       T
// 	Subscribers *SubscriberList[T]
// }

// var _ (Observable[string]) = (*Latch[string])(nil)

// func (x *Latch[T]) Subscribe(done CancellationToken) MessageQueue[T] {
// 	x.Lock.Lock()
// 	defer x.Lock.Unlock()

// 	subscriber := MessageQueue[T](done)

// 	// If value has already arrived, dont bother adding to subscriber list
// 	// Simply emit value and close
// 	if x.Signalled {
// 		subscriber.Post(NewValueMessage[T](x.Value))
// 		subscriber.Post(NewCompleteMessage[T]())
// 		subscriber.Done()
// 		return subscriber
// 	}

// 	x.Subscribers.Subscribe(next, done)
// }

// func (x *Latch[T]) SetValue(value T) {
// 	x.Lock.Lock()
// 	defer x.Lock.Unlock()

// 	for subscriber := range x.Subscribers {
// 		subscriber.Post(NewValueMessage[T](x.Value))
// 		subscriber.Post(NewCompleteMessage[T]())
// 		subscriber.Done()
// 	}

// 	x.Subscribers.Clear()
// }

// func (x *Latch[T]) Unsubscribe(id SubscriberId) {
// 	x.Lock.Lock()
// 	defer x.Lock.Unlock()
// 	x.Subscribers.Unsubscribe(id)
// }
