package rx

type EventType string

const (
	ValueEventType    EventType = "value"
	ErrorEventType    EventType = "error"
	CompleteEventType EventType = "complete"
)

type SubscriberEvent[T any] struct {
	eventType EventType
	value     T
	err       error
}

func NewValueEvent[T any](value T) SubscriberEvent[T] {
	return SubscriberEvent[T]{eventType: ValueEventType, value: value}
}

func NewErrorEvent[T any](err error) SubscriberEvent[T] {
	return SubscriberEvent[T]{eventType: ErrorEventType, err: err}
}

func NewCompleteEvent[T any]() SubscriberEvent[T] {
	return SubscriberEvent[T]{eventType: CompleteEventType}
}

func (x SubscriberEvent[T]) IsValue() (bool, T) {
	return x.eventType == ValueEventType, x.value
}

func (x SubscriberEvent[T]) IsError() (bool, error) {
	return x.eventType == ErrorEventType, x.err
}

func (x SubscriberEvent[T]) IsComplete() bool {
	return x.eventType == CompleteEventType
}
