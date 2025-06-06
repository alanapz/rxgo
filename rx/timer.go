package rx

import (
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func Timer(interval time.Duration) Observable[time.Duration] {
	return NewUnicastObservable(func(args UnicastObserverArgs[time.Duration]) {

		startTime := time.Now()

		for {

			timerChannel := time.After(interval)

			msg := u.SelectReceiveMessage[time.Time]{Policy: u.AbortOnClose}

			if u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectReceive(&timerChannel, &msg)) == u.DoneResult {
				return
			}

			if msg.HasValue && u.Selection(u.SelectDone(args.DownstreamUnsubscribed), u.SelectSend(args.Downstream, msg.Value.Sub(startTime))) == u.DoneResult {
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
