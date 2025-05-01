package rx

import (
	"sync"
	"testing"
	"time"

	"alanpinder.com/rxgo/v2/ux"
)

func TestMerge(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Merge(
		ToAny(Of(1, 2, 3)),
		ToAny(OneShotTimer(5*time.Second)).Pipe(First).Pipe(ConcatMap(func(_ any) Observable[any] { return ToAny(Of(4, 5, 6)) })),
		ToAny(Of(7, 8, 9)),
		ToAny(OneShotTimer(time.Second)).Pipe(First).Pipe(ConcatMap(func(_ any) Observable[any] { return ToAny(Of(10, 11, 12)) })),
	)

	var wg sync.WaitGroup

	done := addTestSubscriber(t, &wg, "s1", source, ux.Of[any](1, 2, 3, 7, 8, 9, 10, 11, 12, 4, 5, 6), ux.Of[error]())

	wg.Wait()
	done()
}
