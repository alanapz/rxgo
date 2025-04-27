package rx

type Void struct{}

func Append[T any](slice *[]T, values ...T) {
	*slice = append(*slice, values...)
}

func Recv1[T any](source1 <-chan T, done <-chan Void) (T, bool, bool) { // msg, isEndOfStream, isDone
	select {
	case _ = <-done:
		return Zero[T](), false, true
	case next, open := <-source1:
		return next, !open, false
	}
}

func RecvWithAdditionalDone[T any, X any](source1 <-chan T, done <-chan Void, done2 <-chan X) (T, bool, bool) { // msg, isEndOfStream, isDone
	select {
	case _ = <-done:
		return Zero[T](), false, true
	case _ = <-done2:
		return Zero[T](), false, true
	case next, open := <-source1:
		return next, !open, false
	}
}

func Send1[T any](value T, dest1 chan<- T, done <-chan Void) bool {
	select {
	case _ = <-done:
		return false
	case dest1 <- value:
		return true
	}
}

func Zero[T any]() T {
	var value T
	return value
}
