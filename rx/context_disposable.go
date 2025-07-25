package rx

import (
	"errors"
	"slices"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type ContextDisposable struct {
	ctx                     *Context
	lock                    *sync.Mutex
	runWorkerOnce           sync.Once
	requestDisposal         chan Void
	closeRequestDisposal    func()
	onDisposalComplete      chan Void
	closeOnDisposalComplete func()
	cleanup                 []func() error
	disposed                bool
	disposalInProgress      bool
	disposalResult          error
}

var _ Disposable = (*ContextDisposable)(nil)

func NewContextDisposable(ctx *Context, lock *sync.Mutex) (*ContextDisposable, error) {

	var requestDisposal chan Void
	var closeRequestDisposal func()

	if err := u.Wrap2(NewChannel[Void](ctx, 1))(&requestDisposal, &closeRequestDisposal); err != nil {
		return nil, err
	}

	var onDisposalComplete chan Void
	var closeOnDisposalComplete func()

	if err := u.Wrap2(NewChannel[Void](ctx, 1))(&onDisposalComplete, &closeOnDisposalComplete); err != nil {
		return nil, err
	}

	disposable := &ContextDisposable{
		ctx:                     ctx,
		lock:                    lock,
		requestDisposal:         requestDisposal,
		closeRequestDisposal:    closeRequestDisposal,
		onDisposalComplete:      onDisposalComplete,
		closeOnDisposalComplete: closeOnDisposalComplete,
	}

	disposable.runWorker()
	return disposable, nil
}

func (x *ContextDisposable) AddCleanup(cleanup func() error) error {

	u.AssertLocked(x.lock)

	if x.disposed {
		return u.ErrDisposed
	}

	if x.disposalInProgress {
		return u.ErrDisposalInProgress
	}

	u.Append(&x.cleanup, cleanup)
	return nil
}

func (x *ContextDisposable) IsDisposed() bool {
	return x.disposed || x.disposalInProgress
}

func (x *ContextDisposable) IsDisposalInProgress() bool {
	return x.disposalInProgress
}

func (x *ContextDisposable) Dispose() {
	x.closeRequestDisposal()
}

func (x *ContextDisposable) OnDisposalComplete() <-chan Void {
	return x.onDisposalComplete
}

func (x *ContextDisposable) runWorker() {
	x.runWorkerOnce.Do(func() {
		GoRun(x.ctx, func() {

			u.Selection2(u.SelectDisposed(x.ctx.Done()), u.SelectDisposed(x.requestDisposal), u.SelectDisposed(x.onDisposalComplete))

			defer u.Lock(x.lock)()
			defer x.closeRequestDisposal()
			defer x.closeOnDisposalComplete()

			if !x.disposed {

				x.disposed = true
				x.disposalInProgress = true

				var cleanupErrors []error

				for cleanup := range slices.Values(x.cleanup) {
					if err := cleanup(); err != nil {
						u.Append(&cleanupErrors, err)
					}
				}

				x.disposalResult = errors.Join(cleanupErrors...)
				x.disposalInProgress = false
			}
		})
	})
}
