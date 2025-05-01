package rx

import (
	"reflect"
	"sync"
	"testing"

	"alanpinder.com/rxgo/v2/ux"
)

func prepareTest(t *testing.T) func() {

	channels, goroutines := GetNumberOfOpenChannels(), GetNumberOfRunningGoRoutines()

	if channels != 0 {
		t.Errorf("Cannot start test with %d channels already open", channels)
	}

	if goroutines != 0 {
		t.Errorf("Cannot start test with %d Go routines already running", goroutines)
	}

	return func() {

		channels, goroutines := GetNumberOfOpenChannels(), GetNumberOfRunningGoRoutines()

		if channels != 0 {
			t.Errorf("%d channels were leaked and should have been closed", channels)
		}

		if goroutines != 0 {
			t.Errorf("%d Go routines were leaked and should have terminated", goroutines)
		}
	}
}

func addTestSubscriber[T any](t *testing.T, wg *sync.WaitGroup, name string, source Observable[T], expectedValues []T, expectedErrors []error) func() {

	values, errors, done := source.Subscribe()

	wg.Add(2)

	GoRun(func() {
		defer wg.Done()

		received := ux.Of[T]()

		for value := range values {
			t.Logf("%s: Received value: '%v'", name, value)
			Append(&received, value)
		}

		t.Logf("%s: Values end-of-stream", name)

		if !reflect.DeepEqual(received, expectedValues) {
			t.Errorf("%s: Unexpected values: '%v' (expected: '%v')", name, received, expectedValues)
		}
	})

	GoRun(func() {
		defer wg.Done()

		received := ux.Of[error]()

		for err := range errors {
			t.Logf("%s: Received error: '%v'", name, err)
			Append(&received, err)
		}

		t.Logf("%s: Errors end-of-stream", name)

		if !reflect.DeepEqual(received, expectedErrors) {
			t.Errorf("%s: Unexpected errors: '%v' (expected: '%v')", name, received, expectedErrors)
		}
	})

	return done
}
