package rx

import (
	"errors"
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func Timer(interval time.Duration) Observable[time.Duration] {
	return NewUnicastObservable(func(ctx *Context, downstream chan<- time.Duration, downstreamUnsubscribed <-chan u.Never) error {

		startTime := time.Now()

		for {

			timerChannel := time.After(interval)

			var timerSignalled bool
			var timerResult time.Time

			if err := u.Selection(ctx,
				u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
				u.SelectReceive(&timerChannel, func(closed bool, value time.Time) error {

					if closed {
						return errors.New("Upstream timer channel closed unexpectedly")
					}

					timerSignalled = true
					timerResult = value
					return nil
				}),
			); err != nil {
				return err
			}

			if timerSignalled {

				if err := u.Selection(ctx,
					u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
					u.SelectSend(downstream, timerResult.Sub(startTime)),
				); err != nil {
					return err
				}
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
