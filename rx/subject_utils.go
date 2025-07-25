package rx

import (
	"slices"

	u "alanpinder.com/rxgo/v2/utils"
)

func sendValuesThenEndOfStreamAsync[T any](ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never, values ...T) {

	GoRun(ctx, func() {

		for value := range slices.Values(values) {

			if err := u.Selection(ctx, u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)), u.SelectSend(downstream, value)); err != nil {
				ctx.Error(err)
				break
			}

		}
	})
}
