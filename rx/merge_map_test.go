package rx

import (
	"sync"
	"testing"
	"time"

	"alanpinder.com/rxgo/v2/ux"
)

func TestMergeMap(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		Of(3, 2, 1, 0),
		MergeMap(func(seconds int) Observable[int] {
			return Pipe(
				OneShotTimer(time.Duration(seconds)*time.Second),
				Map(func(_ time.Time) int {
					return seconds
				}))
		}))

	var wg sync.WaitGroup

	done := addTestSubscriber(t, &wg, "s1", source, ux.Of(0, 1, 2, 3), ux.Of[error]())

	wg.Wait()
	done()
}
