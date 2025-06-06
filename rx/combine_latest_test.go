package rx

import (
	"sync"
	"testing"
	"time"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestCombineLatest(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	type _R = CombineLatestResult[int]

	expected := u.Of(
		u.Of(_R{HasValue: true, Value: 1}, _R{}),
		u.Of(_R{HasValue: true, Value: 2}, _R{}),
		u.Of(_R{HasValue: true, Value: 3}, _R{}),
		u.Of(_R{HasValue: true, Value: 4}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{HasValue: true, Value: 1}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{HasValue: true, Value: 2}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{HasValue: true, Value: 3}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{HasValue: true, Value: 4}),
		u.Of(_R{HasValue: true, Value: 4, EndOfStream: true}, _R{HasValue: true, Value: 4, EndOfStream: true}),
	)

	source := CombineLatest(
		Of(1, 2, 3, 4),
		Pipe(TimerInSeconds(1), Count[int](), Take[int](4)),
	)

	var wg sync.WaitGroup
	var cleanup u.Event

	cleanup.Add(addTestSubscriber(testSubscriberArgs[[]_R]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}

func TestCombineLatestWithReplaySubject(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	type _R = CombineLatestResult[int]

	expected := u.Of(
		u.Of(_R{HasValue: true, Value: 1}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 1}, _R{HasValue: true, Value: 2}, _R{}),
		u.Of(_R{HasValue: true, Value: 1}, _R{HasValue: true, Value: 2}, _R{HasValue: true, Value: 3}),
		u.Of(_R{HasValue: true, Value: 10}, _R{HasValue: true, Value: 2}, _R{HasValue: true, Value: 3}),
		u.Of(_R{HasValue: true, Value: 10}, _R{HasValue: true, Value: 20}, _R{HasValue: true, Value: 3}),
		u.Of(_R{HasValue: true, Value: 10}, _R{HasValue: true, Value: 20}, _R{HasValue: true, Value: 3, EndOfStream: true}),
		u.Of(_R{HasValue: true, Value: 30}, _R{HasValue: true, Value: 20}, _R{HasValue: true, Value: 3, EndOfStream: true}),
		u.Of(_R{HasValue: true, Value: 30}, _R{HasValue: true, Value: 20, EndOfStream: true}, _R{HasValue: true, Value: 3, EndOfStream: true}),
		u.Of(_R{HasValue: true, Value: 30, EndOfStream: true}, _R{HasValue: true, Value: 20, EndOfStream: true}, _R{HasValue: true, Value: 3, EndOfStream: true}),
	)

	subject1 := NewReplaySubject[int](env, 1)
	subject2 := NewReplaySubject[int](env, 1)
	subject3 := NewReplaySubject[int](env, 1)

	source := CombineLatest(subject1, subject2, subject3)

	var wg sync.WaitGroup
	var cleanup u.Event

	u.GoRun(func() {
		env.Error(subject1.Next(1))
		time.Sleep(time.Second)
		env.Error(subject2.Next(2))
		time.Sleep(time.Second)
		env.Error(subject3.Next(3))
		time.Sleep(time.Second)
		env.Error(subject1.Next(10))
		time.Sleep(time.Second)
		env.Error(subject2.Next(20))
		time.Sleep(time.Second)
		env.Error(subject3.EndOfStream())
		time.Sleep(time.Second)
		env.Error(subject1.Next(30))
		time.Sleep(time.Second)
		env.Error(subject2.EndOfStream())
		time.Sleep(time.Second)
		env.Error(subject1.EndOfStream())
	})

	cleanup.Add(addTestSubscriber(testSubscriberArgs[[]_R]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: expected}))

	wg.Wait()
	cleanup.Resolve()
}

func TestCombineLatestWithTakeUntil(t *testing.T) {

	cleanupTest, env := prepareTest(t)
	defer cleanupTest()

	type _R = CombineLatestResult[int]

	expected := u.Of(
		u.Of(_R{HasValue: true, Value: 1}, _R{}),
		u.Of(_R{HasValue: true, Value: 2}, _R{}),
		u.Of(_R{HasValue: true, Value: 3}, _R{}),
		u.Of(_R{HasValue: true, Value: 3, EndOfStream: true}, _R{}),
	)

	subject1 := NewReplaySubject[int](env, 1)
	subject2 := NewReplaySubject[int](env, 1)
	subject3 := NewReplaySubject[Void](env, 1)

	source := Pipe(CombineLatest(subject1, subject2), TakeUntil[[]_R](subject3))

	var wg sync.WaitGroup
	var cleanup u.Event

	u.GoRun(func() {
		env.Error(subject1.Next(1))
		time.Sleep(time.Second)
		env.Error(subject1.Next(2))
		time.Sleep(time.Second)
		env.Error(subject1.Next(3))
		time.Sleep(time.Second)
		env.Error(subject1.EndOfStream())
		time.Sleep(time.Second)
		env.Error(subject3.EndOfStream())
	})

	cleanup.Add(addTestSubscriber(testSubscriberArgs[[]_R]{env: env, name: "s1", t: t, wg: &wg, source: source, expected: expected}))
	wg.Wait()

	env.Error(subject2.EndOfStream())

	cleanup.Resolve()
}
