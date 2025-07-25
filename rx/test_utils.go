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

func prepareTest(t *testing.T) (func(), *Context) {

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
	}, &RxEnvironment{}
}

type testSubscriberArgs[T any] struct {
	ctx      *Context
	name     string
	t        *testing.T
	wg       *sync.WaitGroup
	source   Observable[T]
	expected []T
}

func addTestSubscriber[T any](args testSubscriberArgs[T]) func() {

	t, wg, name, source, expected := args.t, args.wg, args.name, args.source, args.expected

	u.Require(args.env, t, wg, name, source, expected)

	upstream, unsubscribe := source.Subscribe(args.env)

	wg.Add(1)

	u.GoRun(func() {
		defer wg.Done()

		position := 0

		for msg := range upstream {

			t.Logf("%s: (%d) Received: '%v'", name, position+1, msg)

			if position >= len(expected) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected only %d results)", name, position+1, msg, len(expected))
			} else if !reflect.DeepEqual(msg, expected[position]) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected: '%v')", name, position+1, msg, expected[position])
			}

			position++
		}

		t.Logf("%s: End of stream (%d results)", name, position)
	})

	return unsubscribe
}
