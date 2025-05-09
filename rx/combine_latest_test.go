package rx

import (
	"errors"
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestCombineLatest(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	type _R = CombineLatestResult[int]

	source := CombineLatest(
		Of(1, 2, 3, 4),
		Pipe(TimerInSeconds(1), Count[int](), Take[int](4)),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	expected := u.Of(
		u.Of(_R{HasValue: true, Value: 1}, _R{}),
		u.Of(_R{HasValue: true, Value: 2}, _R{}),
		u.Of(_R{HasValue: true, Value: 3}, _R{}),
		u.Of(_R{HasValue: true, Value: 4}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 1}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 2}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 3}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 4}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 4, Complete: true}),
	)

	addTestSubscriber(t, &wg, cleanup, "s1", source, expected, u.Of[error]())

	wg.Wait()
	cleanup.Cleanup()
}

func TestCombineLatestWithError(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	type _R = CombineLatestResult[int]

	err := errors.New("thrown error")

	source := CombineLatest(
		Of(1, 2, 3, 4),
		Pipe(TimerInSeconds(3), First[int](), ConcatMap(func(_ int) Observable[int] { return ThrowError[int](err) })),
		Pipe(TimerInSeconds(2), Count[int](), Take[int](3)),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	expectedValues := u.Of(
		u.Of(_R{HasValue: true, Value: 1}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 2}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 3}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 4}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{}, _R{Value: 1}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 2}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 3}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 4}, _R{}),
		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 4, Complete: true}, _R{}),
	)

	expectedErrors := u.Of(err)

	addTestSubscriber(t, &wg, cleanup, "s1", source, expectedValues, expectedErrors)

	wg.Wait()
	cleanup.Cleanup()
}

// func TestCombineLatestFirst(t *testing.T) {

// 	cleanupTest := prepareTest(t)
// 	defer cleanupTest()

// 	type _R = CombineLatestResult[int]

// 	source := Pipe2(
// 		CombineLatest(
// 			Of(1, 2, 3, 4),
// 			Pipe2(TimerInSeconds(3), Count[int]()),
// 			Pipe2(TimerInSeconds(2), Count[int]()),
// 		),
// 		TakeUntil[[]_R](Pipe2(TimerInSeconds(7), First[int]())),
// 	)

// 	var wg sync.WaitGroup

// 	expected := u.Of(
// 		u.Of(_R{HasValue: true, Value: 1}, _R{}, _R{}),
// 		u.Of(_R{HasValue: true, Value: 2}, _R{}, _R{}),
// 		u.Of(_R{HasValue: true, Value: 3}, _R{}, _R{}),
// 		u.Of(_R{HasValue: true, Value: 4}, _R{}, _R{}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{}, _R{}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{}, _R{HasValue: true, Value: 1}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 1}, _R{HasValue: true, Value: 1}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 1}, _R{HasValue: true, Value: 2}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 2}, _R{HasValue: true, Value: 2}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 2}, _R{HasValue: true, Value: 3}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 2}, _R{HasValue: true, Value: 4}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 3}, _R{HasValue: true, Value: 4}),
// 		u.Of(_R{HasValue: true, Value: 4, Complete: true}, _R{HasValue: true, Value: 3}, _R{HasValue: true, Value: 5}),
// 	)

// 	done := addTestSubscriber(t, &wg, "s1", source, expected, u.Of[error]())

// 	wg.Wait()
// 	done()
// }
