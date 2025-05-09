package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestConcat(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Concat(
		Pipe(TimerInSeconds(3), Take[int](1), Count[int](), ConcatMap(func(_ int) Observable[int] { return Of(1, 2, 3) })),
		Pipe(TimerInSeconds(2), Take[int](1), Count[int](), ConcatMap(func(_ int) Observable[int] { return Of(4, 5, 6) })),
		Pipe(TimerInSeconds(1), Take[int](2), Count[int](), ConcatMap(func(count int) Observable[int] { return Of(7*count, 8*count, 9*count) })),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 14, 16, 18), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
