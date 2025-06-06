package utils

import (
	"fmt"
	"slices"
	"sync"
)

// A CancelFunc tells an operation to abandon its work.
// A CancelFunc does not wait for the work to stop.
// A CancelFunc may be called by multiple goroutines simultaneously.
// After the first call, subsequent calls to a CancelFunc do nothing.
type CancelFunc = func()

type Event struct {
	lock                sync.Mutex
	executionInProgress bool
	executionComplete   bool
	listeners           []*CancelFunc // Pointer so we an easily support removing listeners
	postExecutionHooks  []func()
}

func (x *Event) String() string {
	return fmt.Sprintf("%d listeners, %v %v", len(x.listeners), x.executionInProgress, x.executionComplete)
}

func (x *Event) IsExecutionInProgress() bool {
	return x.executionInProgress
}

func (x *Event) IsResolved() bool {
	return x.executionComplete
}

func (x *Event) Add(listener func()) CancelFunc {

	Assert(listener != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.executionInProgress {
		panic("execution in progress")
	}

	if x.executionComplete {
		panic("execution complete")
	}

	Append(&x.listeners, &listener)

	position := len(x.listeners) - 1

	return func() {
		x.lock.Lock()
		defer x.lock.Unlock()
		x.listeners[position] = nil
	}
}

func (x *Event) Chain(listener func()) CancelFunc {

	Assert(listener != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.executionInProgress {
		panic("execution in progress")
	}

	if x.executionComplete {
		panic("execution complete")
	}

	Append(&x.listeners, &listener)

	position := len(x.listeners) - 1

	return func() {
		x.lock.Lock()
		defer x.lock.Unlock()
		x.listeners[position] = nil
	}
}

func (x *Event) AddPostExecutionHook(postExecutionHook func()) {

	Assert(postExecutionHook != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.executionInProgress {
		panic("execution in progress")
	}

	if x.executionComplete {
		panic("execution complete")
	}

	Append(&x.postExecutionHooks, postExecutionHook)
}

func (x *Event) Resolve() {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.executionInProgress || x.executionComplete {
		return
	}

	x.executionInProgress = true

	for listener := range slices.Values(x.listeners) {
		if listener != nil {
			(*listener)()
		}
	}

	x.executionComplete = true
	x.executionInProgress = false

	for postExecutionHook := range slices.Values(x.postExecutionHooks) {
		postExecutionHook()
	}
}
