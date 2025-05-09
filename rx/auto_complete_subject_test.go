package rx

import (
	"errors"
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestAutoCompleteSubjectAlreadyCompleted(t *testing.T) {

	expectedValue := "hello, world"

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewAutoCompleteSubject[string]()
	subject.PostValue(expectedValue)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s2", subject, u.Of(expectedValue), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}

func TestAutoCompleteSubjectValue(t *testing.T) {

	expectedValue := "hello, world"

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewAutoCompleteSubject[string]()

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s2", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s3", subject, u.Of(expectedValue), u.Of[error]())

	subject.PostValue(expectedValue)

	addTestSubscriber(t, &wg, cleanup, "s4", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s5", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, cleanup, "s6", subject, u.Of(expectedValue), u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}

func TestAutoCompleteSubjectError(t *testing.T) {

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	err := errors.New("new error")

	subject := NewAutoCompleteSubject[string]()

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", subject, u.Of[string](), u.Of(err))
	addTestSubscriber(t, &wg, cleanup, "s2", subject, u.Of[string](), u.Of(err))
	addTestSubscriber(t, &wg, cleanup, "s3", subject, u.Of[string](), u.Of(err))

	subject.PostError(err)

	wg.Wait()
	cleanup.Cleanup()
}

func TestAutoCompleteSubjectReleasesResourcesOnCleanup(t *testing.T) {

	expectedValue := "hello, world"

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewAutoCompleteSubject[string]()

	var wg sync.WaitGroup

	skippedCleanup := &u.Cleanup{} // Use a fake cleanup instead of using u.NewCleanup

	// Note: We deliberately drop reference to unsubscribe !
	addTestSubscriber(t, &wg, skippedCleanup, "s1", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, skippedCleanup, "s2", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, skippedCleanup, "s3", subject, u.Of(expectedValue), u.Of[error]())

	subject.PostValue(expectedValue)

	// Note: We deliberately drop reference to unsubscribe !
	addTestSubscriber(t, &wg, skippedCleanup, "s1", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, skippedCleanup, "s2", subject, u.Of(expectedValue), u.Of[error]())
	addTestSubscriber(t, &wg, skippedCleanup, "s3", subject, u.Of(expectedValue), u.Of[error]())

	wg.Wait()
}
