package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestTake(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	source := Pipe(TimerInSeconds(2), Count[int](), Take[int](4))

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Resolve()
}

func TestTakeWithCustomObservable(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	source := Pipe(NewUnicastObservable(func(args UnicastObserverArgs[int]) {
		t.Logf("New observer subscribed")

		counter := 0

		for {
			counter++
			select {
			case <-args.DownstreamUnsubscribed:
				t.Logf("New observer unsubscribed")
				return
			case args.Downstream <- counter:
				continue
			}
		}

	}), Take[int](4))

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4)}))

	wg.Wait()
	cleanup.Resolve()
}
