package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestReplaySubject(t *testing.T) {

	checkForResourceLeaks, env := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewReplaySubject[int](env, 3)

	env.Error(subject.Next(1))
	env.Error(subject.Next(2))
	env.Error(subject.Next(3))

	var wg sync.WaitGroup
	var cleanup u.Event

	env.Error(subject.Next(4))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: u.Of(2, 3, 4, 5, 6, 7)}))

	env.Error(subject.Next(5))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s2", t: t, wg: &wg, source: subject, expected: u.Of(3, 4, 5, 6, 7)}))

	env.Error(subject.Next(6))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s3", t: t, wg: &wg, source: subject, expected: u.Of(4, 5, 6, 7)}))

	env.Error(subject.Next(7))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s4", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))

	env.Error(subject.EndOfStream())
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s5", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s6", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))

	wg.Wait()
	cleanup.Resolve()
}

func TestReplaySubjectWithLimit(t *testing.T) {

	checkForResourceLeaks, env := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewReplaySubject[int](env, 3)

	env.Error(subject.Next(1))
	env.Error(subject.Next(2))
	env.Error(subject.Next(3))
	env.Error(subject.Next(4))

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: Pipe(subject, Take[int](3)), expected: u.Of(2, 3, 4)}))
	env.Error(subject.EndOfStream())

	wg.Wait()
	cleanup.Resolve()
}
