package utils

import (
	"fmt"
	"reflect"
)

func NewChannel[T any](cleanup *Event, size uint) chan T {

	channel := make(chan T, size)

	onChannelClosed := NewCondition(fmt.Sprintf("Channel: Waiting for channel (type %s, size: %d) to be closed ...", reflect.TypeFor[T]().Name(), size))

	cleanup.Add(func() {
		defer onChannelClosed()
		close(channel)
	})

	return channel
}
