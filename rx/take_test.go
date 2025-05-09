package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestTake(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(TimerInSeconds(2), Count[int](), Take[int](4))

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}

func TestTakeWithCustomObservable(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(NewUnicastObservable(func(values chan<- int, errors chan<- error, unsubscribed <-chan u.Never) {
		t.Logf("New observer subscribed")

		counter := 0

		for {
			counter++
			select {
			case <-unsubscribed:
				t.Logf("New observer unsubscribed")
				return
			case values <- counter:
				continue
			}
		}

	}), Take[int](4))

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
