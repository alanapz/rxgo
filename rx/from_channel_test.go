package rx

import (
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestFromChannel(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	expected := u.Of(0, 1, 2, 3)

	var wg sync.WaitGroup
	var cleanup u.Event

	source := FromChannel(func(unsubscribed <-chan u.Never) <-chan int {

		downstream := make(chan int)

		wg.Add(1)

		u.GoRun(func() {

			defer wg.Done()
			defer close(downstream)

			for value := range slices.Values(expected) {
				select {
				case <-unsubscribed:
					t.Errorf("Unexpected subscribe")
				case downstream <- value:
					continue
				}
			}
		})

		return downstream
	})

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}

func TestFromChannelUnsubscribe(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	expected := u.Of(1, 2, 3, 4)

	var wg sync.WaitGroup
	var cleanup u.Event
	var nextChannelId atomic.Uint32

	channelSupplier := func(unsubscribed <-chan u.Never) <-chan int {

		values := make(chan int)

		channelId := nextChannelId.Add(1)

		wg.Add(1)

		u.GoRun(func() {

			defer wg.Done()
			defer close(values)

			counter := 0

			for {
				counter++

				select {
				case <-unsubscribed:
					t.Logf("%d: Received unsubscribe request", channelId)
					return
				case values <- counter:
					t.Logf("%d: Sent value '%d'", channelId, counter)
					continue
				}
			}
		})

		return values
	}

	source := Pipe(FromChannel(channelSupplier), Take[int](4))

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: source, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: source, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}
