package rx

import (
	"errors"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

var ErrEndOfStream = errors.New("end of stream")
var ErrAborted = errors.New("aborted")

type messageQueue[T any] struct {
	lock             *sync.Mutex
	channel          chan<- T
	aborted          <-chan u.Never
	cleanup          *u.Cleanup
	queue            *u.Queue[u.Optional[T]]
	wakeup           *sync.Cond
	revision         uint64
	result           error
	onWorkerComplete func()
}

func NewMessageQueue[T any](lock *sync.Mutex, channel chan<- T, aborted <-chan u.Never, cleanup *u.Cleanup) *messageQueue[T] {

	queue := &messageQueue[T]{
		lock:             lock,
		channel:          channel,
		aborted:          aborted,
		cleanup:          cleanup,
		queue:            &u.Queue[u.Optional[T]]{},
		wakeup:           sync.NewCond(lock),
		onWorkerComplete: u.NewCondition("Waiting for message queue worker to complete ..."),
	}

	// We always start worker eagerly, so as not to avoid missed cleanups when disposing without events
	u.GoRun(queue.runWorker)
	return queue
}

func (x *messageQueue[T]) Post(value T) error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.queue.Push(u.Optional[T]{HasValue: true, Value: value})
	x.revision++
	x.wakeup.Broadcast()

	return nil
}

func (x *messageQueue[T]) PostComplete() error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.queue.Push(u.Optional[T]{})
	x.revision++
	x.wakeup.Broadcast()

	return nil
}

func (x *messageQueue[T]) DisposeAsync(err error) error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.result = err
	x.revision++
	x.wakeup.Broadcast()

	return nil
}

func (x *messageQueue[T]) Result() error {
	u.AssertLocked(x.lock)
	return x.result
}

func (x *messageQueue[T]) runWorker() {

	u.AssertUnlocked(x.lock)

	var revision uint64
	var err error

	for {
		if err = x.waitForNextRevision(&revision); err != nil {
			break
		}
		if err = x.runWorkerIteration(); err != nil {
			break
		}
	}

	x.disposeWorker(err)
}

func (x *messageQueue[T]) disposeWorker(err error) {

	u.AssertUnlocked(x.lock)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.result == nil {
		x.result = err
	}

	defer x.onWorkerComplete()
	x.cleanup.Do()
}

func (x *messageQueue[T]) waitForNextRevision(current *uint64) error {

	x.lock.Lock()

	for {

		if x.result != nil {
			x.lock.Unlock()
			return x.result
		}

		if *current != x.revision {
			*current = x.revision
			x.lock.Unlock()
			return nil
		}

		x.wakeup.Wait() // Drops then regains lock

		u.AssertLocked(x.lock)
	}
}

// flushQueueIteration returns true if the subscriber should be disposed, or false if it should remain running
func (x *messageQueue[T]) runWorkerIteration() error {

	u.AssertUnlocked(x.lock)

	for {

		var selectItems []u.SelectItem

		if err := x.workerUpdateState(&selectItems); err != nil {
			return err
		}

		if len(selectItems) == 0 {
			return nil
		}

		if u.Selection(selectItems...) == u.DoneResult {
			return ErrAborted
		}
	}
}

func (x *messageQueue[T]) workerUpdateState(selectItems *[]u.SelectItem) error {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.result != nil {
		return x.result
	}

	msg, ok := x.queue.Shift()

	if !ok {
		// Queue empty, nothing to do (for the moment)
		return nil
	}

	if !msg.HasValue {
		// Means we have recieved a complete message - return EndOfStream
		return ErrEndOfStream
	}

	u.Append(selectItems, u.SelectDone(x.aborted), u.SelectSend(x.channel, msg.Value))
	return nil
}
