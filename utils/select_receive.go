package utils

import (
	"fmt"
	"reflect"
	"strings"
)

type SelectReceivePolicy int

const ContinueOnClose SelectReceivePolicy = 0
const AbortOnClose SelectReceivePolicy = 1
const AbortOnMessage SelectReceivePolicy = 2

type SelectReceiveMessage[T any] struct {
	Policy      SelectReceivePolicy
	Selected    bool
	HasValue    bool
	Value       T
	EndOfStream bool
}

var _ fmt.Stringer = (*SelectReceiveMessage[any])(nil)

func (x *SelectReceiveMessage[T]) Reset() {
	x.Selected = false
	x.HasValue = false
	x.Value = Zero[T]()
	x.EndOfStream = false
}

func (x *SelectReceiveMessage[T]) Result() SelectResult {
	if x.EndOfStream && x.Policy&AbortOnClose == AbortOnClose {
		return DoneResult
	}
	if x.HasValue && x.Policy&AbortOnMessage == AbortOnMessage {
		return DoneResult
	}
	return ContinueResult
}

func (x SelectReceiveMessage[T]) String() string {
	var items []string
	if x.Selected {
		Append(&items, "selected")
	}
	if x.HasValue {
		Append(&items, fmt.Sprintf("value=%v", x.Value))
	}
	if x.EndOfStream {
		Append(&items, "end-of-stream")
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ","))
}

func SelectReceive[T any](channel *<-chan T, msg *SelectReceiveMessage[T]) SelectItem {
	return func(x *selectBuilder) SelectResult {

		msg.Reset()

		// If channel is null, means end-of-stream
		if *channel == nil {
			msg.EndOfStream = true
			return msg.Result()
		}

		Append(&x.cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(*channel)})
		Append(&x.handlers, func(received reflect.Value, stillOpen bool) SelectResult {

			msg.Selected = true

			if !stillOpen {
				*channel = nil
				msg.EndOfStream = true
				return msg.Result()
			}

			msg.HasValue = true
			msg.Value = received.Interface().(T)
			msg.EndOfStream = false

			return msg.Result()
		})

		return ContinueResult
	}
}
