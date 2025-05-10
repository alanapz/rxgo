package rx

import (
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestReplaySubject(t *testing.T) {

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewReplaySubject[int](3)

	subject.Next(1)
	subject.Next(2)
	subject.Next(3)

	var wg sync.WaitGroup
	var cleanup u.Event

	subject.Next(4)
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: subject, expected: u.Of(2, 3, 4, 5, 6, 7)}))

	subject.Next(5)
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s2", t: t, wg: &wg, source: subject, expected: u.Of(3, 4, 5, 6, 7)}))

	subject.Next(6)
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s3", t: t, wg: &wg, source: subject, expected: u.Of(4, 5, 6, 7)}))

	subject.Next(7)
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s4", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))

	subject.EndOfStream()
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s5", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))
	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s6", t: t, wg: &wg, source: subject, expected: u.Of(5, 6, 7)}))

	wg.Wait()
	cleanup.Emit()
}

func TestReplaySubjectWithLimit(t *testing.T) {

	checkForResourceLeaks := prepareTest(t)
	defer checkForResourceLeaks()

	subject := NewReplaySubject[int](3)

	subject.Next(1)
	subject.Next(2)
	subject.Next(3)
	subject.Next(4)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[int]{name: "s1", t: t, wg: &wg, source: Pipe(subject, Take[int](3)), expected: u.Of(2, 3, 4)}))
	subject.EndOfStream()

	wg.Wait()
	cleanup.Emit()
}
