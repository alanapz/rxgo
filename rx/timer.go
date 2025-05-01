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

			timerMsg := SelectReceiveMessage[time.Time]{Policy: DoneOnClose}

			timerChannel := time.After(interval)

			if Selection(SelectDone(done), SelectReceive(&timerChannel, &timerMsg)) {
				return
			}

			if timerMsg.HasValue && Selection(SelectDone(done), SelectSend(valuesOut, timerMsg.Value)) == DoneResult {
				return
			}

			if !repeat {
				return
			}
		}
	})
}
