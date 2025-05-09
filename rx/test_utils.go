package rx

import (
	"reflect"
	"slices"
	"sync"
	"testing"

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

			t.Errorf("%d conditions were leaked and should have been resolved", len(conditions))

			for condition := range slices.Values(conditions) {
				t.Errorf("- %s", condition)
			}
		}
	}
}

func addTestSubscriber[T any](t *testing.T, wg *sync.WaitGroup, cleanup *u.Cleanup, name string, source Observable[T], expectedValues []T, expectedErrors []error) {

	values, errors, unsubscribe := source.Subscribe()
	cleanup.Add(unsubscribe)

	wg.Add(2)

	u.GoRun(func() {
		defer wg.Done()

		position := 0

		received := u.Of[T]()

		for value := range values {

			t.Logf("%s: (%d) Received value: '%v'", name, position+1, value)
			u.Append(&received, value)

			if position >= len(expectedValues) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected only %d results)", name, position+1, value, len(expectedValues))
			} else if !reflect.DeepEqual(value, expectedValues[position]) {
				t.Errorf("%s: (%d) Unexpected value: '%v' (expected: '%v')", name, position+1, value, expectedValues[position])
			}

			position++
		}

		t.Logf("%s: Values complete (%d results)", name, position)

		if !reflect.DeepEqual(received, expectedValues) {
			t.Errorf("%s: Unexpected values: '%v' (expected: '%v')", name, received, expectedValues)
		}
	})

	u.GoRun(func() {
		defer wg.Done()

		received := u.Of[error]()

		for err := range errors {
			t.Logf("%s: Received error: '%v'", name, err)
			u.Append(&received, err)
		}

		t.Logf("%s: Errors end-of-stream", name)

		if !reflect.DeepEqual(received, expectedErrors) {
			t.Errorf("%s: Unexpected errors: '%v' (expected: '%v')", name, received, expectedErrors)
		}
	})
}
