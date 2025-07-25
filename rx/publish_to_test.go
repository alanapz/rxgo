package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestPublishToAutoCompleteSubject(t *testing.T) {

	cleanupTest, ctx := prepareTest(t)
	defer cleanupTest()

	// Timers send to Subject1 sends to subject2
	subject := NewAutoCompleteSubject[int](ctx)
	subject2 := NewAutoCompleteSubject[int](ctx)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-1", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-2", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-3", t: t, wg: &wg, source: subject, expected: u.Of(1)}))

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-1", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-2", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-3", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))

	//subject2.AddSource(Pipe(subject, Map(func(value int) int { return value * 2 })), DoNotPropogateEndOfStream)
	subject2.AddSource(subject, DoNotPropogateEndOfStream)
	subject.AddSource(Pipe(TimerInSeconds(1), TapDebug[int](), Count[int]()), DoNotPropogateEndOfStream)

	wg.Wait()
	cleanup.Emit()
}

func TestPublishToAutoCompleteSubject2(t *testing.T) {

	cleanupTest, ctx := prepareTest(t)
	defer cleanupTest()

	// Timers send to Subject1 sends to subject2
	subject := NewAutoCompleteSubject[int](env)
	subject2 := NewAutoCompleteSubject[int](env)
	subject3 := NewAutoCompleteSubject[int](env)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-1", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-2", t: t, wg: &wg, source: subject, expected: u.Of(1)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s1-3", t: t, wg: &wg, source: subject, expected: u.Of(1)}))

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-1", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-2", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s2-3", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s3-1", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s3-2", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{ctx: ctx, name: "s3-3", t: t, wg: &wg, source: subject2, expected: u.Of(2)}))

	subject2.AddSource(Pipe(subject, Map(func(value int) int { return value * 2 })), PropogateEndOfStream)
	subject2.AddSource(Pipe(subject3, Map(func(value int) int { return value * 2 })), PropogateEndOfStream)
	subject.AddSource(Pipe(TimerInSeconds(1), TapDebug[int](), Count[int]()), PropogateEndOfStream)

	wg.Wait()

	cleanup.Emit()
}

// func TestPublishToReplaySubject(t *testing.T) {

// 	cleanupTest, env := prepareTest(t)
// 	defer cleanupTest()

// 	controller := NewReplaySubject[int](env, 3)
// 	subject1 := NewReplaySubject[int](env, 3)
// 	subject2 := NewReplaySubject[int](env, 3)
// 	subject3 := NewReplaySubject[int](env, 3)

// 	var wg sync.WaitGroup
// 	var cleanup u.Event

// 	PublishTo(PublishToArgs[int]{
// 		Environment:          env,
// 		Source:               controller,
// 		Sink:                 subject1,
// 		PropogateEndOfStream: true,
// 	})

// 	PublishTo(PublishToArgs[int]{
// 		Environment: env,
// 		Source:      Pipe(controller, Map(func(value int) int { return value * 10 })),
// 		Sink:        subject2,
// 	})

// 	PublishTo(PublishToArgs[int]{
// 		Environment:          env,
// 		Source:               subject1,
// 		Sink:                 subject3,
// 		PropogateEndOfStream: true,
// 	})

// 	PublishTo(PublishToArgs[int]{
// 		Environment: env,
// 		Source:      subject2,
// 		Sink:        subject3,
// 	})

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{
// 		name:       "s1",
// 		t:          t,
// 		wg:         &wg,
// 		source:     subject3,
// 		expected:   u.Of(1, 2),
// 		outOfOrder: true,
// 	}))

// 	env.Error(controller.Next(1))
// 	env.Error(controller.Next(2))
// 	env.Error(controller.EndOfStream())

// 	wg.Wait()

// 	env.Error(subject3.EndOfStream())

// 	cleanup.Resolve()
// }
