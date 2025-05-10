package utils

import "sync"

type Queue[T any] struct {
	lock  sync.Mutex
	items []T
}

func (x *Queue[T]) Clear() {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.items = []T{}
}

func (x *Queue[T]) Shift() (T, bool) {

	x.lock.Lock()
	defer x.lock.Unlock()

	if len(x.items) == 0 {
		return Zero[T](), false
	}

	first := x.items[0]
	x.items = x.items[1:len(x.items)]
	return first, true
}

func (x *Queue[T]) Push(values ...T) {

	x.lock.Lock()
	defer x.lock.Unlock()

	Append(&x.items, values...)
}
