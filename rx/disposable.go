package rx

import u "alanpinder.com/rxgo/v2/utils"

type Disposable interface {
	IsDisposed() bool
	IsDisposalInProgress() bool
	Dispose()
	OnDisposalComplete() <-chan u.Void
}
