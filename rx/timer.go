package rx

import "time"

func Timer(interval time.Duration) Observable[time.Time] {
	return NewUnicastObservable(func(observer chan<- Message[time.Time], done <-chan Void) {

		for {

			timer, endOfStream, isDone := Recv1(time.After(interval), done)

			if endOfStream || isDone {
				return
			}

			if !Send1(NewValue(timer), observer, done) {
				return
			}

		}
	})
}
