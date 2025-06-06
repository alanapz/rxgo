package rx

import (
	"fmt"
	"reflect"

	u "alanpinder.com/rxgo/v2/utils"
)

type RxEnvironment struct {
	CustomExecutor     func(func())
	CustomErrorHandler func(error)
}

func (x *RxEnvironment) Execute(fn func()) {
	u.Coalesce(x.CustomExecutor, u.GoRun)(fn)
}

func (x *RxEnvironment) Error(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func NewChannel[T any](env *RxEnvironment, bufferCapacity uint) (chan T, *u.Event) {

	channel := make(chan T, bufferCapacity)

	onChannelClosed := u.NewCondition(fmt.Sprintf("Channel: Waiting for channel (type %s, size: %d) to be closed ...", reflect.TypeFor[T]().Name(), bufferCapacity))

	var channelCloseEvent u.Event

	channelCloseEvent.Add(func() {
		defer onChannelClosed()
		close(channel)
	})

	return channel, &channelCloseEvent
}
