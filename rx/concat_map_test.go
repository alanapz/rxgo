package rx

import (
	"sync"
	"testing"
	"time"

	"alanpinder.com/rxgo/v2/ux"
)

func TestConcatMap(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		Of(2, 1, 0),
		ConcatMap(func(seconds int) Observable[int] {
			return Pipe(
				OneShotTimer(time.Duration(seconds)*time.Second),
				Map(func(_ time.Time) int {
					return seconds
				}))
		}))

	wg := &sync.WaitGroup{}

	done := addTestSubscriber(t, wg, "s1", source, ux.Of(2, 1, 0), ux.Of[error]())

	wg.Wait()
	done()
}
