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
		Pipe(TimerInSeconds(4), Take[int](1), Count[int](), ConcatMap(func(_ int) Observable[int] { return Of(1, 2, 3) })),
		Pipe(TimerInSeconds(3), Take[int](1), Count[int](), ConcatMap(func(_ int) Observable[int] { return Of(4, 5, 6) })),
		Pipe(TimerInSeconds(1), Take[int](2), Count[int](), ConcatMap(func(count int) Observable[int] { return Of(7*count, 8*count, 9*count) })),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(7, 8, 9, 14, 16, 18, 4, 5, 6, 1, 2, 3)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(7, 8, 9, 14, 16, 18, 4, 5, 6, 1, 2, 3)}))

	wg.Wait()
	cleanup.Emit()
}

func TestMergeRace(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	// Subject1 wins
	func() {

		subject1 := NewReplaySubject[int](1)
		subject2 := NewReplaySubject[int](1)

		source := Pipe(Merge(subject1, subject2), First[int]())

		var wg sync.WaitGroup
		var cleanup u.Event

		wg.Add(1)

		u.GoRun(func() {
			defer wg.Done()

			subject1.Next(1)
			time.Sleep(time.Second)

			subject2.Next(2)
			time.Sleep(time.Second)

			subject1.EndOfStream()
			subject2.EndOfStream()
		})

		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1)}))
		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1)}))

		wg.Wait()
		cleanup.Emit()
	}()

	// Subject2 wins
	func() {

		subject1 := NewReplaySubject[int](1)
		subject2 := NewReplaySubject[int](1)

		source := Pipe(Merge(subject1, subject2), First[int]())

		var wg sync.WaitGroup
		var cleanup u.Event

		wg.Add(1)

		u.GoRun(func() {
			defer wg.Done()

			subject2.Next(2)
			time.Sleep(time.Second)

			subject1.Next(1)
			time.Sleep(time.Second)

			subject1.EndOfStream()
			subject2.EndOfStream()
		})

		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(2)}))
		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(2)}))

		wg.Wait()
		cleanup.Emit()
	}()

}
