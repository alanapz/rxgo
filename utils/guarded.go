package utils

import "sync"

type GuardedBy[T any] struct {
	lock  *sync.Mutex
	value T
}

func NewGuardedBy[T any](lock *sync.Mutex, value T) *GuardedBy[T] {
	return &GuardedBy[T]{lock: lock, value: value}
}

func (x *GuardedBy[T]) Get() T {
	AssertLocked(x.lock)
	return x.value
}

func (x *GuardedBy[T]) Set(value T) {
	AssertLocked(x.lock)
	x.value = value
}
