package rx

type NewUnicastObserver[T any] = func(values chan<- T, errors chan<- error, done <-chan Never)

type UnicastObservable[T any] struct {
	onNewObserver NewUnicastObserver[T]
}

var _ Observable[any] = (*UnicastObservable[any])(nil)

func NewUnicastObservable[T any](onNewObserver NewUnicastObserver[T]) *UnicastObservable[T] {
	return &UnicastObservable[T]{
		onNewObserver: onNewObserver,
	}
}

func (x *UnicastObservable[T]) Subscribe() (<-chan T, <-chan error, func()) {

	values, cleanupValues := NewChannel[T](0)
	errors, cleanupErrors := NewChannel[error](0)
	done, cleanupDone := NewChannel[Never](0)

	go func() {
		defer cleanupValues()
		defer cleanupErrors()
		x.onNewObserver(values, errors, done)
	}()

	return values, errors, cleanupDone
}

func (x *UnicastObservable[T]) Pipe(operator OperatorFunction[T, T]) Observable[T] {
	return operator(x)
}
