package utils

import (
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"
)

type ItemId uint64

type ItemList[T any] struct {
	lock       sync.Mutex
	nextItemId ItemId
	items      map[ItemId]T
}

func (x *ItemList[T]) Add(item T) func() {

	defer Lock(&x.lock)()

	x.nextItemId++

	itemId := x.nextItemId
	x.items[itemId] = item

	return sync.OnceFunc(func() {
		defer Lock(&x.lock)()
		delete(x.items, itemId)
	})
}

func (x *ItemList[T]) Items() iter.Seq[T] {

	defer Lock(&x.lock)()

	if len(x.items) == 0 {
		return nil
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
