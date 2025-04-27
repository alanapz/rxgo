package rx

import "fmt"

type MessageType string

var ValueMessage MessageType = "value"
var ErrorMessage MessageType = "error"

type Message[T any] struct {
	Type  MessageType
	Value T
	Error error
}

func NewValue[T any](value T) Message[T] {
	return Message[T]{Type: ValueMessage, Value: value}
}

func NewError[T any](err error) Message[T] {
	return Message[T]{Type: ErrorMessage, Error: err}
}

func MapMessage[T any, U any](msg Message[T], valueProjection func(T) U) Message[U] {
	if msg.IsValue() {
		return Message[U]{Type: ValueMessage, Value: valueProjection(msg.Value), Error: msg.Error}
	}
	return Message[U]{Type: msg.Type, Value: Zero[U](), Error: msg.Error}
}

func (x Message[T]) IsValue() bool {
	return x.Type == ValueMessage
}

func (x Message[T]) IsError() bool {
	return x.Type == ErrorMessage
}

func (x Message[T]) String() string {
	if x.IsValue() {
		return fmt.Sprintf("value %v", x.Value)
	}
	if x.IsError() {
		return fmt.Sprintf("error %v", x.Error)
	}
	return fmt.Sprintf("%s", x.Type)
}
