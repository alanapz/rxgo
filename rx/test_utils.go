package rx

import (
	"reflect"
	"slices"
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func prepareTest(t *testing.T) func() {

	conditions := getConditions()

	if len(conditions) != 0 {

		t.Errorf("Cannot start test with %d conditions unresolved", len(conditions))

		for condition := range slices.Values(conditions) {
			t.Errorf("- %s", condition)
		}
	}

	return func() {

		conditions := getConditions()

		if len(conditions) != 0 {

			t.Errorf("%d conditions were leaked and should have been resolved", len(conditions))

			for condition := range slices.Values(conditions) {
				t.Errorf("- %s", condition)
			}
		}
	}
}

func addTestSubscriber[T any](t *testing.T, wg *sync.WaitGroup, name string, source Observable[T], expectedValues []T, expectedErrors []error) func() {

	values, errors, done := source.Subscribe()

	wg.Add(2)

	GoRun(func() {
		defer wg.Done()

		received := u.Of[T]()

		for value := range values {
			t.Logf("%s: Received value: '%v'", name, value)
			u.Append(&received, value)
		}

		if !reflect.DeepEqual(received, expectedValues) {
			t.Errorf("%s: Unexpected values: '%v' (expected: '%v')", name, received, expectedValues)
		}
	})

	GoRun(func() {
		defer wg.Done()

		received := u.Of[error]()

		for err := range errors {
			t.Logf("%s: Received error: '%v'", name, err)
			u.Append(&received, err)
		}

		if !reflect.DeepEqual(received, expectedErrors) {
			t.Errorf("%s: Unexpected errors: '%v' (expected: '%v')", name, received, expectedErrors)
		}
	})

	return done
}
