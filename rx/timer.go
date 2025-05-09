package rx

import (
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func Timer(interval time.Duration) Observable[time.Duration] {
	return NewUnicastObservable(func(valuesOut chan<- time.Duration, errorsOut chan<- error, unsubscribed <-chan u.Never) {

		startTime := time.Now()

		for {

			timerChannel := time.After(interval)

			timerMsg := u.SelectReceiveMessage[time.Time]{Policy: u.AbortOnClose}

			if u.Selection(u.SelectDone(unsubscribed), u.SelectReceive(&timerChannel, &timerMsg)) == u.DoneResult {
				return
			}

			if timerMsg.HasValue && u.Selection(u.SelectDone(unsubscribed), u.SelectSend(valuesOut, timerMsg.Value.Sub(startTime))) == u.DoneResult {
				return
			}
		}
	})
}

func TimerInSeconds(seconds int) Observable[int] {
	return Pipe2(
		Timer(time.Duration(seconds)*time.Second),
		Map(func(duration time.Duration) int {
			return int(duration.Seconds())
		}))
}
