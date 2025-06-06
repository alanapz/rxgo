package rx

import (
	"slices"
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestAutoCompleteSubjectAlreadyCompleted(t *testing.T) {

	checkForResourceLeaks, env := prepareTest(t)
	defer checkForResourceLeaks()

	expected := u.Of("hello, world")

	subject := NewAutoCompleteSubject[string](env)

	for expectedValue := range slices.Values(expected) {
		env.Error(subject.Next(expectedValue))
	}

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s2", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s3", t: t, wg: &wg, source: subject, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}

func TestAutoCompleteSubject(t *testing.T) {

	checkForResourceLeaks, env := prepareTest(t)
	defer checkForResourceLeaks()

	expected := u.Of("hello, world")

	subject := NewAutoCompleteSubject[string](env)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s2", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s3", t: t, wg: &wg, source: subject, expected: expected}))

	for expectedValue := range slices.Values(expected) {
		env.Error(subject.Next(expectedValue))
	}

	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s4", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s5", t: t, wg: &wg, source: subject, expected: expected}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s6", t: t, wg: &wg, source: subject, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}

func TestAutoCompleteSubjectReleasesResourcesOnCleanup(t *testing.T) {

	checkForResourceLeaks, env := prepareTest(t)
	defer checkForResourceLeaks()

	expected := u.Of("hello, world")

	subject := NewAutoCompleteSubject[string](env)

	var wg sync.WaitGroup

	_ = addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: expected}) // Note dont attach to cleanup
	_ = addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s2", t: t, wg: &wg, source: subject, expected: expected})
	_ = addTestSubscriber(testSubscriberArgs[string]{env: env, name: "s3", t: t, wg: &wg, source: subject, expected: expected})

	for expectedValue := range slices.Values(expected) {
		env.Error(subject.Next(expectedValue))
	}

	wg.Wait()
}
