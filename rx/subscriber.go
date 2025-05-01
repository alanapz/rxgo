package rx

import (
	"sync"
)

type SubscriberId = uint64

type Subscriber[T any] struct {
	lock          *sync.Mutex
	values        chan<- T
	errors        chan<- error
	done          <-chan Never
	cleanup       func()
	queue         *Queue[SubscriberEvent[T]]
	disposed      bool
	activeWorkers uint64
}

func NewSubscriber[T any](lock *sync.Mutex, values chan<- T, errors chan<- error, done <-chan Never, cleanup func()) *Subscriber[T] {
	return &Subscriber[T]{
		lock:    lock,
		values:  values,
		errors:  errors,
		done:    done,
		cleanup: cleanup,
		queue:   NewQueue[SubscriberEvent[T]](),
	}
}

func (x *Subscriber[T]) PostEvent(event SubscriberEvent[T]) bool {

	AssertLocked(x.lock)

	if x.disposed {
		return true
	}

	x.queue.Push(event)

	GoRun(x.flushQueue)

	return false
}

func (x *Subscriber[T]) flushQueue() {

	workerNecessary, cleanup := x.isWorkerNecessary()
	defer cleanup()

	if !workerNecessary {
		return
	}

	if x.flushQueueIteration() {
		x.dispose()
	}
}

func (x *Subscriber[T]) dispose() {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		panic("subscriber already disposed")
	}

	x.disposed = true
	x.cleanup()
}

// flushQueueIteration returns true if the subscriber should be disposed, or false if it should remain running
func (x *Subscriber[T]) flushQueueIteration() bool {

	for {

		var selectItems []SelectItem

		if x.buildNextSelection(&selectItems) {
			return true
		}

		if len(selectItems) == 0 {
			return false
		}

		if Selection(selectItems...) {
			return true
		}
	}
}

func (x *Subscriber[T]) buildNextSelection(selectItems *[]SelectItem) bool {

	x.lock.Lock()
	defer x.lock.Unlock()

	event, ok := x.queue.PopFirst()

	if !ok {
		// Queue empty, nothing to do (for the moment)
		return false
	}

	if event.IsComplete() {
		return true
	}

	isValue, value := event.IsValue()
	isError, err := event.IsError()

	Append(selectItems, SelectDone(x.done))

	if isValue {
		Append(selectItems, SelectSend(x.values, value))
	}

	if isError {
		Append(selectItems, SelectSend(x.errors, err))
	}

	return false
}

func (x *Subscriber[T]) isWorkerNecessary() (bool, func()) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.disposed {
		return false, NoOp
	}

	if x.activeWorkers > 0 {
		return false, NoOp
	}

	x.activeWorkers++

	return true, func() {

		x.lock.Lock()
		defer x.lock.Unlock()

		x.activeWorkers--
	}
}
