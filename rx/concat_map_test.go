package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestConcatMap(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	source := Pipe(
		Of(3, 2, 1, 0),
		ConcatMap(func(seconds int) Observable[int] {
			return Pipe(TimerInSeconds(seconds), MapTo[int](seconds), First[int]())
		}))

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(3, 2, 1, 0)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(3, 2, 1, 0)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(3, 2, 1, 0)}))

	wg.Wait()
	cleanup.Resolve()
}
