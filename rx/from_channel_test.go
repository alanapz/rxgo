package rx

import (
	"slices"
	"sync"
	"sync/atomic"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestFromChannel(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	expectedValues := u.Of(0, 1, 2, 3)
	expectedErrors := u.Of[error]()

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	source := FromChannel(func(unsubscribed <-chan u.Never) (<-chan int, <-chan error) {

		values := make(chan int)
		errors := make(chan error)

		wg.Add(2)

		u.GoRun(func() {

			defer wg.Done()
			defer close(values)

			for value := range slices.Values(expectedValues) {
				select {
				case <-unsubscribed:
					t.Errorf("Unexpected subscribe")
				case values <- value:
					continue
				}
			}
		})

		u.GoRun(func() {
			defer wg.Done()
			defer close(errors)
		})

		return values, errors
	})

	addTestSubscriber(t, &wg, cleanup, "s1", source, expectedValues, expectedErrors)
	addTestSubscriber(t, &wg, cleanup, "s2", source, expectedValues, expectedErrors)
	addTestSubscriber(t, &wg, cleanup, "s2", source, expectedValues, expectedErrors)

	wg.Wait()
	cleanup.Cleanup()
}

func TestFromChannelUnsubscribe(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	var nextChannelId atomic.Uint32

	channelSupplier := func(unsubscribed <-chan u.Never) (<-chan int, <-chan error) {

		values := make(chan int)
		errors := make(chan error)

		wg.Add(2)

		channelId := nextChannelId.Add(1)

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

		u.GoRun(func() {
			defer wg.Done()
			defer close(errors)
		})

		return values, errors
	}

	source := Pipe(FromChannel(channelSupplier), Take[int](4))

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(1, 2, 3, 4), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s2", source, u.Of(1, 2, 3, 4), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s3", source, u.Of(1, 2, 3, 4), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}
