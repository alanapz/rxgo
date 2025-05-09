package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestMergeMap(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe2(
		Of(3, 2, 1, 0),
		MergeMap(func(seconds int) Observable[int] {
			return Pipe(TimerInSeconds(seconds), First[int]())
		}))

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(0, 1, 2, 3), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
