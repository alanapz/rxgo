package rx

// import (
// 	"sync"
// 	"testing"
// 	"time"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// func TestBufferedSubject(t *testing.T) {

// 	cleanupTest, ctx := prepareTest(t)
// 	defer cleanupTest()

// 	bufferReleased := NewNotificationSubject(ctx)
// 	defer bufferReleased.Dispose()

// 	source := Pipe(
// 		TimerInSeconds(1),
// 		Tap(func(value int) {
// 			t.Logf("Before buffer released, value: %v, signalled: %v", value, bufferReleased.IsSignalled())

// 			if bufferReleased.IsSignalled() {
// 				t.Error("Was not expecting buffer to be released before all values passed through first tap")
// 			}

// 			if value == 10 {
// 				go func() {
// 					time.Sleep(1 * time.Second)
// 					t.Logf("Releasing buffer ...")
// 					ctx.Assert(bufferReleased.Signal())
// 				}()
// 			}
// 		}),
// 		Take[int](10),
// 		BufferUntil[int](bufferReleased),
// 		Tap(func(value int) {
// 			t.Logf("After buffer released, value: %v, signalled: %v", value, bufferReleased.IsSignalled())

// 			if !bufferReleased.IsSignalled() {
// 				t.Error("Was expecting buffer to be released before values passing through second tap")
// 			}
// 		}))

// 	var wg sync.WaitGroup
// 	var cleanup u.Event

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1", t: t, wg: &wg, source: source, expected: u.Of(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)}))

// 	wg.Wait()
// 	cleanup.Emit()
// }
