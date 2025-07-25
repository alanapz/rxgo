package rx

import "errors"

var ErrDownstreamUnsubscribed = errors.New("Downstream unsubscribed")

type Observable[T any] interface {
	Subscribe(ctx *Context) (<-chan T, func(), error)
}
