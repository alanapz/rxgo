package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestTakeUntilWithTimer(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		TimerInSeconds(2),
		Count[int](),
		TakeUntil[int](MapTo[int](Void{})(Pipe(TimerInSeconds(3), First[int]()))),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2)}))

	wg.Wait()
	cleanup.Resolve()
}

func TestTakeUntilWithSubjectAfter(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	func() {

		notifier1 := NewNotificationSubject(env)
		defer notifier1.Dispose()

		notifier2 := NewNotificationSubject(env)
		defer notifier2.Dispose()

		source := Pipe(
			TimerInSeconds(1),
			Count[int](),
			TakeUntil[int](notifier1, notifier2),
			Tap(func(value int) {
				if value == 4 {
					env.Error(notifier2.Signal())
				}
			}),
		)

		var wg sync.WaitGroup
		var cleanup u.Event

		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
		cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

		wg.Wait()
		cleanup.Resolve()
	}()
}

func TestTakeUntilWithSubjectBefore(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[Void](env)

	source := Pipe(
		TimerInSeconds(1),
		Count[int](),
		Tap(func(value int) {
			if value == 4 {
				env.Error(subject.Next(Void{}))
			}
		}),
		TakeUntil[int](subject),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Resolve()
}
