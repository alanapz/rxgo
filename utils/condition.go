package utils

import (
	"fmt"
	"maps"
	"runtime/debug"
	"slices"
	"sync"
)

var conditionLock sync.Mutex
var nextConditionId uint64
var conditions = map[uint64]condition{}

type condition struct {
	message string
	stack   string
}

var _ fmt.Stringer = (*condition)(nil)

func NewCondition(message string) func() {

	conditionLock.Lock()
	defer conditionLock.Unlock()

	nextConditionId++

	conditionId := nextConditionId
	conditions[conditionId] = condition{message: message, stack: string(debug.Stack())}

	return func() {

		conditionLock.Lock()
		defer conditionLock.Unlock()

		if _, ok := conditions[conditionId]; !ok {
			panic(fmt.Sprintf("condition already resolved: %d", conditionId))
		}

		delete(conditions, conditionId)
	}
}

func (x condition) String() string {
	return x.message + x.stack
}

func GetConditions() []condition {

	conditionLock.Lock()
	defer conditionLock.Unlock()

	return slices.Collect(maps.Values(conditions))
}
