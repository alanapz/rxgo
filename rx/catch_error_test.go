package rx

import (
	"errors"
	"sync"
	"testing"

	u "alanpinder.com/rxgo/v2/utils"
)

func TestCatchError(t *testing.T) {

	cleanupTest := prepareTest(t)
	defer cleanupTest()

	err1 := errors.New("error#1")
	err1bis := errors.New("error#1bis")
	err2 := errors.New("error#2")
	err3 := errors.New("error#3")

	source := Pipe(
		Of(0, 1, 2, 3, 4),
		ConcatMap(func(value int) Observable[int] {
			if value == 1 {
				return ThrowError[int](err1)
			}
			if value == 2 {
				return ThrowError[int](err2)
			}
			if value == 3 {
				return ThrowError[int](err3)
			}
			return Of(value)
		}),
		CatchError(func(err error) Observable[int] {
			if err == err1 {
				return ThrowError[int](err1bis)
			}
			if err == err3 {
				return Of(3)
			}
			return ThrowError[int](err)
		}),
	)

	var wg sync.WaitGroup

	cleanup := u.NewCleanup(t.Name())

	addTestSubscriber(t, &wg, cleanup, "s1", source, u.Of(0, 3, 4), u.Of(err1bis, err2))

	wg.Wait()
	cleanup.Cleanup()
}
