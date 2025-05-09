package utils

import (
	"sync"
)

func RunOnce(fn func()) func() {

	var once sync.Once

	onResolved := NewCondition("RunOnce: Waiting for cleanup function to be called ...")

	return func() {
		once.Do(func() {
			defer onResolved()
			fn()
		})
	}
}

func RunAtMostOnce(fn func()) func() {

	var lock sync.Mutex
	var signalled bool

	onResolved := NewCondition("RunAtMostOnce: Waiting for cleanup function to be called (at most once) ...")

	return func() {
		defer onResolved()

		lock.Lock()
		defer lock.Unlock()

		if signalled {
			panic("Cleanup function already called")
		}

		signalled = true
		fn()
	}
}
