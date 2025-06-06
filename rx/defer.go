package rx

func Defer[T any](provider func() Observable[T]) Observable[T] {
	return NewUnicastObservable(func(args UnicastObserverArgs[T]) {
		drainObservable(drainObservableArgs[T]{
			Environment:            args.Environment,
			Source:                 provider(),
			Downstream:             args.Downstream,
			DownstreamUnsubscribed: args.DownstreamUnsubscribed,
		})
	})
}
