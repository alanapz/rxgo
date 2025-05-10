package utils

import (
	"fmt"
	"slices"
	"sync"
)

type Event struct {
	lock                sync.Mutex
	executionInProgress bool
	executionComplete   bool
	listeners           []*func() // Pointer so we an easily support removing listeners
	postExecutionHooks  []func()
}

func (x *Event) String() string {
	return fmt.Sprintf("%d listeners, %v %v", len(x.listeners), x.executionInProgress, x.executionComplete)
}

func (x *Event) IsExecutionInProgress() bool {
	return x.executionInProgress
}

func (x *Event) IsExecutionComplete() bool {
	return x.executionComplete
}

func (x *Event) Add(listener func()) func() {

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

func (x *Event) Emit() {

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
