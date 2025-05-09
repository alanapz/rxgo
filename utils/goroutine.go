package utils

import (
	"fmt"
)

func GoRun(fn func()) {

	onComplete := NewCondition(fmt.Sprintf("Waiting for goroutine %p to exit", fn))

	go func() {
		defer onComplete()
		fn()
	}()
}
