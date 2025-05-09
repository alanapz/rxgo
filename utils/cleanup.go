package utils

import (
	"fmt"
	"slices"
	"sync"
)

type Cleanup struct {
	lock              sync.Mutex
	label             string
	cleanupInProgress bool
	cleanupComplete   bool
	cleaners          []func()
	postCleanupHooks  []func()
	onCleanedUp       func()
}

func NewCleanup(label string) *Cleanup {
	return &Cleanup{
		label:       label,
		onCleanedUp: NewCondition(fmt.Sprintf("Cleanup: Waiting for %s to run", Ternary(label != "", fmt.Sprintf("%s cleanup", label), "cleanup"))),
	}
}

func (x *Cleanup) IsInProgress() bool {
	return x.cleanupInProgress
}

func (x *Cleanup) IsComplete() bool {
	return x.cleanupComplete
}

func (x *Cleanup) Add(cleaner func()) {

	Assert(cleaner != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.cleanupInProgress {
		panic("cleanup in progress")
	}

	if x.cleanupComplete {
		panic("cleanup already complete")
	}

	Append(&x.cleaners, cleaner)
}

func (x *Cleanup) AddPostCleanupHook(postCleanupHook func()) {

	Assert(postCleanupHook != nil)

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.cleanupInProgress {
		panic("cleanup in progress")
	}

	if x.cleanupComplete {
		panic("cleanup already complete")
	}

	Append(&x.postCleanupHooks, postCleanupHook)
}

func (x *Cleanup) Do() {
	x.Cleanup()
}

func (x *Cleanup) Cleanup() {

	x.lock.Lock()
	defer x.lock.Unlock()

	if x.cleanupInProgress || x.cleanupComplete {
		return
	}

	x.cleanupInProgress = true

	for cleaner := range slices.Values(x.cleaners) {
		cleaner()
	}

	x.cleanupInProgress = false
	x.cleanupComplete = true

	defer x.onCleanedUp()

	for postCleanupHook := range slices.Values(x.postCleanupHooks) {
		postCleanupHook()
	}
}
