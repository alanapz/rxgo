package rx

import (
	"maps"
	"slices"
	"sync"
)

var _conditionLock sync.Mutex
var _nextConditionId uint64
var _conditions = map[uint64]string{}

func AddCondition(message string) func() {

	_conditionLock.Lock()
	defer _conditionLock.Unlock()

	_nextConditionId++

	conditionId := _nextConditionId
	_conditions[conditionId] = message

	return func() {
		_conditionLock.Lock()
		defer _conditionLock.Unlock()
		delete(_conditions, conditionId)
	}
}

func getConditions() []string {
	_conditionLock.Lock()
	defer _conditionLock.Unlock()
	return slices.Collect(maps.Values(_conditions))
}
