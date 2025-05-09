package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestTakeUntilWithTimer(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		TimerInSeconds(2),
		Count[int](),
		TakeUntil[int](Pipe(TimerInSeconds(5), First[int]())),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}

func TestTakeUntilWithSubject(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[bool]()

	source := Pipe(
		TimerInSeconds(1),
		Count[int](),
		TakeUntil[int](subject),
		Tap(func(value int) {
			if value == 4 {
				subject.PostValue(true)
			}
		}, nil),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
