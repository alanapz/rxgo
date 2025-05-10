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
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Emit()
}

func TestTakeWithCustomObservable(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	source := Pipe(NewUnicastObservable(func(downstream chan<- int, unsubscribed <-chan u.Never) {
		t.Logf("New observer subscribed")

		counter := 0

		for {
			counter++
			select {
			case <-unsubscribed:
				t.Logf("New observer unsubscribed")
				return
			case downstream <- counter:
				continue
			}
		}

	}), Take[int](4))

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Emit()
}
