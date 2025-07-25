package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func Of[T any](values ...T) Observable[T] {
	return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

		for value := range slices.Values(values) {

			if err := u.Selection(ctx, u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)), u.SelectSend(downstream, value)); err != nil {
				return err
			}

		}

		return nil
	})
}
