package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestPublishToAutoCompleteSubject(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	subject := NewAutoCompleteSubject[int](env)
	chainedSubject := NewAutoCompleteSubject[int](env)

	var wg sync.WaitGroup
	var cleanup u.Event

	PublishTo(PublishToArgs[int]{
		Environment: env,
		Source:      Pipe(TimerInSeconds(1), TapDebug[int](), Count[int]()),
		Sink:        subject,
	})

	PublishTo(PublishToArgs[int]{
		Environment: env,
		Source:      Pipe(subject, Map(func(value int) int { return value * 2 })),
		Sink:        chainedSubject,
	})

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "s1", t: t, wg: &wg, source: subject, expected: u.Of(1)}))

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "c1", t: t, wg: &wg, source: chainedSubject, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "c2", t: t, wg: &wg, source: chainedSubject, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{env: env, name: "c3", t: t, wg: &wg, source: chainedSubject, expected: u.Of(2)}))

	wg.Wait()
	cleanup.Resolve()
}

func TestPublishToReplaySubject(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	controller := NewReplaySubject[int](env, 3)
	subject1 := NewReplaySubject[int](env, 3)
	subject2 := NewReplaySubject[int](env, 3)
	subject3 := NewReplaySubject[int](env, 3)

	var wg sync.WaitGroup
	var cleanup u.Event

	PublishTo(PublishToArgs[int]{
		Environment:          env,
		Source:               controller,
		Sink:                 subject1,
		PropogateEndOfStream: true,
	})

	PublishTo(PublishToArgs[int]{
		Environment: env,
		Source:      Pipe(controller, Map(func(value int) int { return value * 10 })),
		Sink:        subject2,
	})

	PublishTo(PublishToArgs[int]{
		Environment:          env,
		Source:               subject1,
		Sink:                 subject3,
		PropogateEndOfStream: true,
	})

	PublishTo(PublishToArgs[int]{
		Source: subject2,
		Sink:   subject3,
	})

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{
		name:       "s1",
		t:          t,
		wg:         &wg,
		source:     subject3,
		expected:   u.Of(1, 2),
		outOfOrder: true,
	}))

	env.Error(controller.Next(1))
	env.Error(controller.Next(2))
	env.Error(controller.EndOfStream())

	wg.Wait()

	env.Error(subject3.EndOfStream())

	cleanup.Resolve()
}
