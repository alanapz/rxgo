package rx

import (
	"reflect"
	"runtime"
	"slices"
	"sync"
	"testing"
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func prepareTest(t *testing.T) func() {

	conditions := u.GetConditions()

	if len(conditions) != 0 {

		t.Errorf("Cannot start test with %d conditions unresolved", len(conditions))

		for condition := range slices.Values(conditions) {
			t.Errorf("- %s", condition)
		}
	}

	return func() {

		conditions := u.GetConditions()

		if len(conditions) != 0 {
			runtime.Gosched()
			time.Sleep(time.Second * 2)
		}

		conditions = u.GetConditions()

		if len(conditions) != 0 {

			t.Errorf("%d conditions were leaked and should have been resolved", len(conditions))

			for condition := range slices.Values(conditions) {
				t.Errorf("- %s", condition)
			}
		}
	}
}

type testSubscriberArgs[T any] struct {
	name       string
	t          *testing.T
	wg         *sync.WaitGroup
	source     Observable[T]
	expected   []T
	outOfOrder bool
}

func addTestSubscriber[T any](args testSubscriberArgs[T]) func() {

	t, wg, name, source, expected, outOfOrder := args.t, args.wg, args.name, args.source, args.expected, args.outOfOrder

	u.Require(t, wg, name, source, expected)

	upstream, unsubscribe := source.Subscribe()

	wg.Add(1)

	u.GoRun(func() {
		defer wg.Done()

		position := 0

		received := u.Of[T]()

		for msg := range upstream {

			t.Logf("%s: (%d) Received: '%v'", name, position+1, msg)
			u.Append(&received, msg)

			if position >= len(expected) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected only %d results)", name, position+1, msg, len(expected))
			} else if !outOfOrder && !reflect.DeepEqual(msg, expected[position]) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected: '%v')", name, position+1, msg, expected[position])
			}

			position++
		}

		t.Logf("%s: End of stream (%d results)", name, position)

		if !reflect.DeepEqual(received, expected) {
			t.Errorf("%s: Unexpected values: '%v' (expected: '%v')", name, received, expected)
		}
	})

	return unsubscribe
}
