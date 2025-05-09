package utils

import (
	"fmt"
	"reflect"
)

func NewChannel2[T any](size uint) (chan T, func()) {

	channel := make(chan T, size)

	onChannelClosed := NewCondition(fmt.Sprintf("Channel: Waiting for channel (type %s, size: %d) to be closed ...", reflect.TypeFor[T]().Name(), size))

	return channel, RunAtMostOnce(func() {
		defer onChannelClosed()
		close(channel)
	})
}

func NewChannel[T any](cleanup *Cleanup, size uint) chan T {

	channel := make(chan T, size)

	onChannelClosed := NewCondition(fmt.Sprintf("Channel: Waiting for channel (type %s, size: %d) to be closed ...", reflect.TypeFor[T]().Name(), size))

	cleanup.Add(func() {
		defer onChannelClosed()
		close(channel)
	})

	return channel
}
