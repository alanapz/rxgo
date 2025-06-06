package rx

import (
	"errors"
	"fmt"
	"slices"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

var ErrEndOfStream = errors.New("end of stream")
var ErrAborted = errors.New("aborted")
var ErrDisposed = errors.New("disposed")

type messageValue[T any] struct {
	value       T
	endOfStream bool
}

var _ fmt.Stringer = (*messageValue[any])(nil)

func (x messageValue[T]) String() string {
	if x.endOfStream {
		return "{end-of-stream}"
	}
	return fmt.Sprintf("{%v}", x.value)
}

type messageQueue[T any] struct {
	lock             *sync.Mutex
	channel          chan<- T
	aborted          <-chan u.Never
	cleanup          *u.Event
	queue            *u.Queue[messageValue[T]]
	wakeup           *sync.Cond
	revision         uint64
	result           error
	onWorkerComplete func()
}

func NewMessageQueue[T any](lock *sync.Mutex, channel chan<- T, aborted <-chan u.Never, cleanup *u.Event, initial []messageValue[T]) *messageQueue[T] {

	queue := &messageQueue[T]{
		lock:             lock,
		channel:          channel,
		aborted:          aborted,
		cleanup:          cleanup,
		queue:            &u.Queue[messageValue[T]]{},
		wakeup:           sync.NewCond(lock),
		onWorkerComplete: u.NewCondition("Waiting for message queue worker to complete ..."),
	}

	if len(initial) > 0 {
		queue.queue.Push(initial...)
		queue.revision++
	}

	// We always start worker eagerly, so as not to avoid missed cleanups when disposing without events
	u.GoRun(queue.runWorker)
	return queue
}

func (x *messageQueue[T]) Next(values ...T) error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	for value := range slices.Values(values) {
		x.queue.Push(messageValue[T]{value: value})
		x.revision++
	}

	x.wakeup.Broadcast()

	return nil
}

func (x *messageQueue[T]) EndOfStream() error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.queue.Push(messageValue[T]{endOfStream: true})
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
	x.cleanup.Resolve()
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

	// println(fmt.Printf("workerUpdateState msg=%v ok=%v", msg, ok))

	if !ok {
		// Queue empty, nothing to do (for the moment)
		return nil
	}

	if msg.endOfStream {
		// Means we have recieved a complete message - return EndOfStream
		return ErrEndOfStream
	}

	u.Append(selectItems, u.SelectDone(x.aborted), u.SelectSend(x.channel, msg.value))
	return nil
}
