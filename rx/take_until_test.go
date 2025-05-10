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
		TakeUntil[int](Pipe(TimerInSeconds(3), First[int]())),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))

	wg.Wait()
	cleanup.Emit()
}

func TestTakeUntilWithSubjectAfter(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[struct{}]()

	source := Pipe(
		TimerInSeconds(1),
		Count[int](),
		TakeUntil[int](subject),
		Tap(func(value int) {
			if value == 4 {
				subject.Next(struct{}{})
			}
		}),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Emit()
}

func TestTakeUntilWithSubjectBefore(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[struct{}]()

	source := Pipe(
		TimerInSeconds(1),
		Count[int](),
		Tap(func(value int) {
			if value == 4 {
				subject.Next(struct{}{})
			}
		}),
		TakeUntil[int](subject),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Emit()
}
