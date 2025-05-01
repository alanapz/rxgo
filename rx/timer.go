package rx

import (
	"time"
)

func OneShotTimer(interval time.Duration) Observable[time.Time] {
	return newTimer(interval, false)
}

func RepeatableTimer(interval time.Duration) Observable[time.Time] {
	return newTimer(interval, true)
}

func newTimer(interval time.Duration, repeat bool) Observable[time.Time] {
	return NewUnicastObservable(func(valuesOut chan<- time.Time, errorsOut chan<- error, done <-chan Never) {
		for {

			var isValue bool
			var value time.Time

			if Selection(SelectDone(done), SelectMustReceive(time.After(interval), &isValue, &value)) {
				return
			}

			if isValue && Selection(SelectDone(done), SelectSend(valuesOut, value)) {
				return
			}

			if !repeat {
				return
			}
		}
	})
}
