package rx

import (
	"sync"
	"testing"
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestMerge(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Merge(
		ToAny(ConcatMap(func(_ time.Time) Observable[int] { return Of(10, 11, 12) })(OneShotTimer(4*time.Second))),
		ToAny(ConcatMap(func(_ time.Time) Observable[int] { return Of(7, 8, 9) })(OneShotTimer(3*time.Second))),
		ToAny(ConcatMap(func(_ time.Time) Observable[int] { return Of(4, 5, 6) })(OneShotTimer(2*time.Second))),
		ToAny(ConcatMap(func(_ time.Time) Observable[int] { return Of(1, 2, 3) })(OneShotTimer(time.Second))),
	)

	var wg sync.WaitGroup

	done := addTestSubscriber(t, &wg, "s1", source, u.Of[any](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12), u.Of[error]())

	wg.Wait()
	done()
}
