package utils

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type ListenerFunc func()
type CancelFunc func()

type EventState string

const ExecutionInProgress EventState = "executionInProgress"
const RunningPostExecutionHooks EventState = "runningPostExecutionHooks"
const ExecutionComplete EventState = "complete"

type Event struct {
	lock               sync.Mutex
	state              EventState
	nextListenerId     uint
	listeners          map[uint]ListenerFunc
	postExecutionHooks map[uint]ListenerFunc
}

func (x *Event) String() string {
	return fmt.Sprintf("%d listeners, state=%v", len(x.listeners), x.state)
}

func (x *Event) IsInProgress() bool {
	return x.state == ExecutionInProgress || x.state == RunningPostExecutionHooks
}

func (x *Event) IsComplete() bool {
	return x.state == ExecutionComplete
}

func (x *Event) Add(listener ListenerFunc) CancelFunc {

	Assert(listener != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.IsInProgress() {
		panic("execution in progress")
	}

	if x.IsComplete() {
		panic("execution complete")
	}

	listenerId := x.nextListenerId
	x.nextListenerId++

	x.listeners[listenerId] = listener

	return func() {

		x.lock.Lock()
		defer x.lock.Unlock()

		delete(x.listeners, listenerId)
	}
}

func (x *Event) AddPostExecutionHook(postExecutionHook ListenerFunc) CancelFunc {

	Assert(postExecutionHook != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.IsInProgress() {
		panic("execution in progress")
	}

	if x.IsComplete() {
		panic("execution complete")
	}

	listenerId := x.nextListenerId
	x.nextListenerId++

	x.postExecutionHooks[listenerId] = postExecutionHook

	return func() {

		x.lock.Lock()
		defer x.lock.Unlock()

		delete(x.postExecutionHooks, listenerId)
	}
}

func (x *Event) Emit() {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.IsInProgress() || x.IsComplete() {
		return
	}

	x.state = ExecutionInProgress

	if len(x.listeners) > 0 {
		listenerIds := slices.Collect(maps.Keys(x.listeners))
		slices.Sort(listenerIds)
		for listenerId := range slices.Values(listenerIds) {
			x.listeners[listenerId]()
		}
	}

	if len(x.postExecutionHooks) > 0 {

		x.state = RunningPostExecutionHooks

		listenerIds := slices.Collect(maps.Keys(x.postExecutionHooks))
		slices.Sort(listenerIds)
		for listenerId := range slices.Values(listenerIds) {
			x.postExecutionHooks[listenerId]()
		}
	}

	x.state = ExecutionComplete
}
