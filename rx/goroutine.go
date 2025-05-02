package rx

import (
	"fmt"
)

func GoRun(fn func()) {

	clearCondition := AddCondition(fmt.Sprintf("Waiting for goroutine %p to exit", fn))

	go func() {
		defer clearCondition()
		fn()
	}()
}
