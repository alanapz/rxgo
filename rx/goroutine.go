package rx

import (
	"sync/atomic"
)

var _runningGoRoutines atomic.Int32

func GoRun(fn func()) {

	_runningGoRoutines.Add(1)

	go func() {
		defer _runningGoRoutines.Add(-1)
		fn()
	}()
}

func GetNumberOfRunningGoRoutines() int32 {
	return _runningGoRoutines.Load()
}
