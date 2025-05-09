package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestConcatMap(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		Of(3, 2, 1, 0),
		ConcatMap(func(seconds int) Observable[int] {
			return Pipe(TimerInSeconds(seconds), MapTo[int](seconds), First[int]())
		}))

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(3, 2, 1, 0), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
