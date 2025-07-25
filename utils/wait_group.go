package utils

import (
	"errors"
	"sync"
)

type WaitGroup struct {
	wg      sync.WaitGroup
	errLock sync.Mutex
	errors  []error
}

func (x *WaitGroup) Add(delta int) {
	x.wg.Add(delta)
}

func (x *WaitGroup) Done() {
	x.wg.Done()
}

func (x *WaitGroup) Wait() {
	x.wg.Wait()
}

func (x *WaitGroup) Error(err error) {
	defer Lock(&x.errLock)()
	Append(&x.errors, err)
}

func (x *WaitGroup) GetError() error {
	defer Lock(&x.errLock)()
	return errors.Join(x.errors...)
}
