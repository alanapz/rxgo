package rx

import (
	"sync"
	"sync/atomic"
)

var _openChannels atomic.Int32

func NewChannel[T any](size uint) (chan T, func()) {

	channel := make(chan T, size)
	_openChannels.Add(1)

	var cleanup sync.Once

	return channel, func() {
		cleanup.Do(func() {
			close(channel)
			_openChannels.Add(-1)
		})
	}
}

func GetNumberOfOpenChannels() int32 {
	return _openChannels.Load()
}
