package rx

import (
	"fmt"
	"sync"
)

func NewChannel[T any](size uint) (chan T, func()) {

	channel := make(chan T, size)

	clearCondition := AddCondition(fmt.Sprintf("Waiting for channel %p to be closed", channel))

	var cleanup sync.Once

	return channel, func() {
		cleanup.Do(func() {
			defer clearCondition()
			close(channel)
		})
	}
}
