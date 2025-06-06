package rx

type Disposable interface {
	IsDisposed() bool
	Dispose()
}
