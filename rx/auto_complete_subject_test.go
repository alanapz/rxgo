package rx

// import (
// 	"slices"
// 	"sync"
// 	"testing"

// 	u "alanpinder.com/rxgo/v2/utils"
// )

// func TestAutoCompleteSubjectAlreadyCompleted(t *testing.T) {

// 	checkForResourceLeaks, ctx := prepareTest(t)
// 	defer checkForResourceLeaks()

// 	expected := u.Of("hello, world")

// 	subject := NewAutoCompleteSubject[string](ctx)

// 	for expectedValue := range slices.Values(expected) {
// 		ctx.Assert(subject.Next(expectedValue))
// 	}

// 	var wg sync.WaitGroup

// 	cleanup, emitCleanup := ctx.NewEvent()

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s1", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s2", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s3", t: t, wg: &wg, source: subject, expected: expected}))

// 	wg.Wait()
// 	emitCleanup()
// }

// func TestAutoCompleteSubject(t *testing.T) {

// 	checkForResourceLeaks, ctx := prepareTest(t)
// 	defer checkForResourceLeaks()

// 	expected := u.Of("hello, world")

// 	subject := NewAutoCompleteSubject[string](ctx)

// 	var wg sync.WaitGroup
// 	var cleanup u.Event

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s1", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s2", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s3", t: t, wg: &wg, source: subject, expected: expected}))

// 	for expectedValue := range slices.Values(expected) {
// 		ctx.Assert(subject.Next(expectedValue))
// 	}

// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s4", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s5", t: t, wg: &wg, source: subject, expected: expected}))
// 	cleanup.Add(addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s6", t: t, wg: &wg, source: subject, expected: expected}))

// 	wg.Wait()
// 	cleanup.Emit()
// }

// func TestAutoCompleteSubjectReleasesResourcesOnCleanup(t *testing.T) {

// 	checkForResourceLeaks, ctx := prepareTest(t)
// 	defer checkForResourceLeaks()

// 	expected := u.Of("hello, world")

// 	subject := NewAutoCompleteSubject[string](ctx)

// 	var wg sync.WaitGroup

// 	_ = addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s1", t: t, wg: &wg, source: subject, expected: expected}) // Note dont attach to cleanup
// 	_ = addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s2", t: t, wg: &wg, source: subject, expected: expected})
// 	_ = addTestSubscriber(testSubscriberArgs[string]{ctx: ctx, name: "s3", t: t, wg: &wg, source: subject, expected: expected})

// 	for expectedValue := range slices.Values(expected) {
// 		ctx.Assert(subject.Next(expectedValue))
// 	}

// 	wg.Wait()
// }
