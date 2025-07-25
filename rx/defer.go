package rx

import u "alanpinder.com/rxgo/v2/utils"

func Defer[T any](provider func() (Observable[T], error)) Observable[T] {
	return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

		var observable Observable[T]

		if err := u.Wrap(provider())(&observable); err != nil {
			return nil
		}

		return drainObservable(drainObservableArgs[T]{
			Context:                ctx,
			Source:                 observable,
			Downstream:             downstream,
			DownstreamUnsubscribed: downstreamUnsubscribed,
		})
	})
}
