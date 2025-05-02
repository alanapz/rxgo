package rx

import (
	"errors"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

var ErrDisposed = errors.New("Already disposed")
var ErrEndOfStream = errors.New("End of stream")

type messageQueue[T any] struct {
	lock      *sync.Mutex
	channel   chan<- T
	done      <-chan u.Never
	cleanup   func(error)
	queue     *u.Queue[u.Optional[T]]
	wakeup    *sync.Cond
	revision  uint64
	result    error
	cleanedUp sync.Once
}

func NewMessageQueue[T any](lock *sync.Mutex, channel chan<- T, done <-chan u.Never, cleanup func(error)) *messageQueue[T] {

	queue := &messageQueue[T]{
		lock:    lock,
		channel: channel,
		done:    done,
		cleanup: cleanup,
		queue:   u.NewQueue[u.Optional[T]](),
		wakeup:  sync.NewCond(lock),
	}

	// We always start worker eagerly, so as not to avoid missec cleanups when disposing without events
	GoRun(queue.runWorker)
	return queue
}

func (x *messageQueue[T]) PostValue(value T) error {
	return x.doPostValue(u.NewOptional(value))
}

func (x *messageQueue[T]) PostComplete() error {
	return x.doPostValue(u.EmptyOptional[T]())
}

func (x *messageQueue[T]) doPostValue(value u.Optional[T]) error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.queue.Push(value)
	x.revision++
	x.wakeup.Broadcast()

	return nil
}

func (x *messageQueue[T]) Dispose(err error) error {

	u.AssertLocked(x.lock)

	if x.result != nil {
		return x.result
	}

	x.result = err
	x.revision++
	x.wakeup.Broadcast()

	return nil
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

	x.disposeWorker(err)()
}

func (x *messageQueue[T]) disposeWorker(err error) func() {

	u.Assert(err != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.result == nil {
		x.result = err
	}

	return func() {
		u.AssertUnlocked(x.lock)
		x.cleanedUp.Do(func() {
			x.cleanup(x.result)
		})
	}
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

		if u.Selection(selectItems...) {
			return ErrDisposed
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

	u.Append(selectItems, u.SelectDone(x.done), u.SelectSend(x.channel, msg.Value))
	return nil
}
