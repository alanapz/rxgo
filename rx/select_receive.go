package rx

import (
	"fmt"
	"reflect"
)

type SelectReceivePolicy bool

const ContinueOnClose SelectReceivePolicy = false
const DoneOnClose SelectReceivePolicy = true

type SelectReceiveMessage[T any] struct {
	Policy      SelectReceivePolicy
	HasValue    bool
	Value       T
	EndOfStream bool
}

func (x *SelectReceiveMessage[T]) Reset() {
	x.HasValue = false
	x.Value = Zero[T]()
	x.EndOfStream = false
}

func (x SelectReceivePolicy) endOfStreamResult() SelectResult {
	if x == ContinueOnClose {
		return ContinueResult
	}
	if x == DoneOnClose {
		return DoneResult
	}
	panic(fmt.Sprintf("unexpected policy: %v", x))
}

func SelectReceive[T any](channel *<-chan T, msg *SelectReceiveMessage[T]) SelectItem {
	return func(x *selectBuilder) SelectResult {

		msg.Reset()

		// If channel is null, means end-of-stream
		if *channel == nil {
			msg.EndOfStream = true
			return msg.Policy.endOfStreamResult()
		}

		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*channel)})
		Append(&x.handlers, func(received reflect.Value, stillOpen bool) SelectResult {

			if !stillOpen {
				*channel = nil
				msg.EndOfStream = true
				return msg.Policy.endOfStreamResult()
			}

			msg.HasValue = true
			msg.Value = received.Interface().(T)
			msg.EndOfStream = false

			return ContinueResult
		})

		return ContinueResult
	}
}
