package rx

// import (
// 	"sync"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// var stateOpen = 1
// var stateOpen = 1

// type subscriber[T any] struct {
// 	lock            *sync.Mutex
// 	values          *MessageQueue[T]
// 	cleanupValues   func()
// 	valuesCleanedUp sync.Once
// 	errors          *MessageQueue[error]
// 	cleanupErrors   func()
// 	done            <-chan u.Never
// 	result          error
// }

// func NewSubscriber[T any](lock *sync.Mutex, values chan<- T, errors chan<- error, done <-chan Never, valuesCleanup func()) *Subscriber[T] {

// 	var valuesCleanedUp sync.Once

// 	valuesDisposed := func(err error) {
// 		valuesCleanedUp.Do(func() {
// 			valuesCleanup()
// 		})
// 	}

// 	errorDisposed := func(err error) {
// 		valuesCleanedUp.Do(func() {
// 			valuesCleanup()
// 		})
// 	}

// 	return &Subscriber[T]{
// 		lock:    lock,
// 		values:  NewMessageQueue(lock, values, done),
// 		errors:  errors,
// 		done:    done,
// 		cleanup: cleanup,
// 		queue:   NewQueue[SubscriberEvent[T]](),
// 	}
// }

// func (x *subscriber[T]) PostValue(value T) error {

// 	u.AssertLocked(x.lock)

// 	if x.result != nil {
// 		return x.result
// 	}

// 	return x.values.PostValue(value)
// }

// func (x *subscriber[T]) PostError(value T) error {

// 	u.AssertLocked(x.lock)

// 	if x.result != nil {
// 		return x.result
// 	}

// 	return x.values.PostValue(value)
// }

// func (x *Subscriber[T]) PostComplete() error {

// 	AssertLocked(x.lock)

// 	if x.result != nil {
// 		return x.result
// 	}

// 	return x.values.PostValue(value)
// 	x.queue.Push(event)

// 	GoRun(x.flushQueue)

// 	return false
// }

// func (x *Subscriber[T]) flushQueue() {

// 	workerNecessary, cleanup := x.isWorkerNecessary()
// 	defer cleanup()

// 	if !workerNecessary {
// 		return
// 	}

// 	if x.flushQueueIteration() {
// 		x.dispose()
// 	}
// }

// func (x *Subscriber[T]) dispose() {

// 	x.lock.Lock()
// 	defer x.lock.Unlock()

// 	if x.disposed {
// 		panic("subscriber already disposed")
// 	}

// 	x.disposed = true
// 	x.cleanup()
// }

// // flushQueueIteration returns true if the subscriber should be disposed, or false if it should remain running
// func (x *Subscriber[T]) flushQueueIteration() bool {

// 	for {

// 		var selectItems []SelectItem

// 		if x.buildNextSelection(&selectItems) {
// 			return true
// 		}

// 		if len(selectItems) == 0 {
// 			return false
// 		}

// 		if Selection(selectItems...) {
// 			return true
// 		}
// 	}
// }

// func (x *Subscriber[T]) buildNextSelection(selectItems *[]SelectItem) bool {

// 	x.lock.Lock()
// 	defer x.lock.Unlock()

// 	event, ok := x.queue.PopFirst()

// 	if !ok {
// 		// Queue empty, nothing to do (for the moment)
// 		return false
// 	}

// 	if event.IsComplete() {
// 		return true
// 	}

// 	isValue, value := event.IsValue()
// 	isError, err := event.IsError()

// 	Append(selectItems, SelectDone(x.done))

// 	if isValue {
// 		Append(selectItems, SelectSend(x.values, value))
// 	}

// 	if isError {
// 		Append(selectItems, SelectSend(x.errors, err))
// 	}

// 	return false
// }

// func (x *Subscriber[T]) isWorkerNecessary() (bool, func()) {

// 	x.lock.Lock()
// 	defer x.lock.Unlock()

// 	if x.disposed {
// 		return false, NoOp
// 	}

// 	if x.activeWorkers > 0 {
// 		return false, NoOp
// 	}

// 	x.activeWorkers++

// 	return true, func() {

// 		x.lock.Lock()
// 		defer x.lock.Unlock()

// 		x.activeWorkers--
// 	}
// }
