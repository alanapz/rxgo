package utils

import (
	"fmt"
	"maps"
	"slices"
	"sync"
)

type ListenerFunc func()

type EventState = string

const ExecutionInProgress EventState = "executionInProgress"
const RunningPostExecutionHooks EventState = "runningPostExecutionHooks"
const ExecutionComplete EventState = "complete"

type EventEmitter struct {
	lock               sync.Mutex
	state              EventState
	nextListenerId     uint
	listeners          map[uint]ListenerFunc
	postExecutionHooks map[uint]ListenerFunc
}

func (x *EventEmitter) NewEventEmitter() *EventEmitter {
	return x.state == ExecutionInProgress || x.state == RunningPostExecutionHooks
}

func (x *EventEmitter) IsExecutionInProgress() bool {
	return x.state == ExecutionInProgress || x.state == RunningPostExecutionHooks
}

func (x *EventEmitter) IsExecutionComplete() bool {
	return x.state == ExecutionComplete
}

func (x *EventEmitter) AddListener(listener ListenerFunc) {

	Assert(listener != nil)

	defer Lock(&x.lock)()

	if x.IsExecutionInProgress() {
		panic("execution in progress")
	}

	if x.IsComplete() {
		panic("execution complete")
	}

	listenerId := x.nextListenerId

	x.nextListenerId++
	x.listeners[listenerId] = listener
}

func (x *EventEmitter) AddPostExecutionHook(postExecutionHook ListenerFunc) {

	Assert(postExecutionHook != nil)

	defer Lock(&x.lock)()

	if x.IsExecutionInProgress() {
		panic("execution in progress")
	}

	if x.IsComplete() {
		panic("execution complete")
	}

	listenerId := x.nextListenerId

	x.nextListenerId++
	x.postExecutionHooks[listenerId] = postExecutionHook
}

func (x *EventEmitter) Emit() {

	defer Lock(&x.lock)()

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

func (x *EventEmitter) String() string {
	return fmt.Sprintf("%d listeners, state=%v", len(x.listeners), x.state)
}
