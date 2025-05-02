package rx

import (
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func OneShotTimer(interval time.Duration) Observable[time.Time] {
	return newTimer(interval, false)
}

func RepeatableTimer(interval time.Duration) Observable[time.Time] {
	return newTimer(interval, true)
}

func newTimer(interval time.Duration, repeat bool) Observable[time.Time] {
	return NewUnicastObservable(func(valuesOut chan<- time.Time, errorsOut chan<- error, done <-chan u.Never) {
		for {

			timerChannel := time.After(interval)

			timerMsg := u.SelectReceiveMessage[time.Time]{Policy: u.DoneOnClose}

			if u.Selection(u.SelectDone(done), u.SelectReceive(&timerChannel, &timerMsg)) {
				return
			}

			if timerMsg.HasValue && u.Selection(u.SelectDone(done), u.SelectSend(valuesOut, timerMsg.Value)) == u.DoneResult {
				return
			}

			if !repeat {
				return
			}
		}
	})
}
