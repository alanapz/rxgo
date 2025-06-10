package rx

import (
	"fmt"
	"reflect"

	u "alanpinder.com/rxgo/v2/utils"
)

type RxEnvironment struct {
	parent       *RxEnvironment
	executor     func(func())
	errorHandler func(error)
	cleanup      u.Event
}

func (x *RxEnvironment) NewEnvironment() *RxEnvironment {
	return &RxEnvironment{
		executor:     u.GoRun,
		errorHandler: u.Panic,
	}
}

func (x *RxEnvironment) AddCleanup(listener func()) {
	x.cleanup.Add(listener)
}

func (x *RxEnvironment) AddTryCleanup(listener func() error) {
	x.cleanup.Add(listener)
}

func (x *RxEnvironment) Dispose() {
	x.cleanup.Emit()
}

func (x *RxEnvironment) Execute(fn func()) {

	u.Assert(fn != nil)

	if x.executor != nil {
		x.executor(fn)
		return
	}

	x.parent.Execute(fn)
}

func (x *RxEnvironment) Error(err error) {

	if err == nil {
		return
	}

	if x.errorHandler != nil {
		x.errorHandler(err)
		return
	}

	x.parent.Error(err)
}

func (x *RxEnvironment) WithExecutor(executor func(func())) *RxEnvironment {
	u.Assert(executor != nil)
	return &RxEnvironment{parent: x, executor: executor}
}

func (x *RxEnvironment) WithErrorHandler(errorHandler func(error)) *RxEnvironment {
	u.Assert(errorHandler != nil)
	return &RxEnvironment{parent: x, errorHandler: errorHandler}
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
