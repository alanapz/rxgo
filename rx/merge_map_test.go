package rx

// import (
// 	"sync"
// 	"testing"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// func TestMergeMap(t *testing.T) {

// 	cleanupTest, env := prepareTest(t)
// 	defer cleanupTest()

// 	source := Pipe2(
// 		Of(3, 2, 1, 0),
// 		MergeMap(func(seconds int) Observable[int] {
// 			return Pipe(TimerInSeconds(seconds), First[int]())
// 		}))

// 	var wg sync.WaitGroup

// 	cleanup := u.Event{}

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(0, 1, 2, 3)}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: u.Of(0, 1, 2, 3)}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: u.Of(0, 1, 2, 3)}))

// 	wg.Wait()
// 	cleanup.Emit()
// }
