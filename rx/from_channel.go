package rx

import (
	u "alanpinder.com/rxgo/v2/utils"
)

func FromChannel[T any](channelSupplier func() <-chan T) Observable[T] {
	return FromChannelWithError(func(c <-chan u.Never) (<-chan T, error) {
		return channelSupplier(), nil
	})
}

func FromChannelWithError[T any](channelSupplier func(<-chan u.Never) (<-chan T, error)) Observable[T] {
	return NewUnicastObservable(func(ctx *Context, downstream chan<- T, downstreamUnsubscribed <-chan u.Never) error {

		var upstream <-chan T

		if err := u.Wrap(channelSupplier(downstreamUnsubscribed))(&upstream); err != nil {
			return err
		}

		for {

			if upstream == nil {
				return nil
			}

			var channelHasValue bool
			var channelValue T

			if err := u.Selection(ctx,
				u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)),
				u.SelectReceive(&upstream, func(closed bool, value T) error {

					if closed {
						return nil
					}

					channelHasValue = true
					channelValue = value
					return nil

				})); err != nil {
				return err
			}

			if channelHasValue {

				if err := u.Selection(ctx, u.SelectDone(downstreamUnsubscribed, u.Val(ErrDownstreamUnsubscribed)), u.SelectSend(downstream, channelValue)); err != nil {
					return err
				}

			}
		}
	})
}
