package rx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"

	u "alanpinder.com/rxgo/v2/utils"
)

type CleanupFunction = func() error
type Executor = func(func())

type Context struct {
	lock           sync.Mutex
	disposed       bool
	notifyDisposed chan Void
	cleanup        u.ItemList[CleanupFunction]
	resources      u.ItemList[ContextResource]
	executor       Executor
	errorHandler   func(error)
}

type ContextResource struct {
	message string
	stack   string
}

func (x *Context) IsDisposed() bool {
	return x.disposed
}

func (x *Context) Dispose() error {

	x.lock.Lock()

	// Do nothing if we are already disposed
	if x.disposed {
		x.lock.Unlock()
		return nil
	}

	x.disposed = true
	x.lock.Unlock()

	var errs []error

	for cleanup := range x.cleanup.Items() {
		cleanup()
	}

	return errors.Join(errs...)
}

func (x *Context) Err() error {

	defer u.Lock(&x.lock)()

	if x.disposed {
		return ErrDisposed
	}

	return nil
}

func (x *Context) Done() <-chan struct{} {
	return nil
}

func (x *Context) NewContext() *Context {
	return &Context{
		executor: func(fn func()) {
			go fn()
		},
		errorHandler: func(err error) {
			panic(err)
		},
	}
}

func (x *Context) AddCleanup(cleanup CleanupFunction) error {

	defer u.Lock(&x.lock)()

	if x.disposed {
		return ErrDisposed
	}

	x.cleanup.Add(cleanup)
	return nil
}

func (x *Context) Assert(err error) {
	x.Error(err)
}

func (x *Context) Error(err error) {
	if err != nil && !errors.Is(err, context.Canceled) {
		x.errorHandler(err)
	}
}

func (x *Context) clone() *Context {
	return &Context{
		executor:     x.executor,
		errorHandler: x.errorHandler,
	}
}

func WithExecutor(ctx *Context, executor func(func())) *Context {
	u.Assert(executor != nil)
	newContext := ctx.clone()
	newContext.executor = executor
	return newContext
}

func WithErrorHandler(ctx *Context, errorHandler func(error)) *Context {
	u.Assert(errorHandler != nil)
	newContext := ctx.clone()
	newContext.errorHandler = errorHandler
	return newContext
}

func GoRun(ctx *Context, fn func()) error {
	return GoRunWithLabel(ctx, "(anonymous)", fn)
}

func GoRunWithLabel(ctx *Context, label string, fn func()) error {

	var resourceCleanup func()

	if err := u.Get(NewResource(ctx, fmt.Sprintf("Waiting for goroutine '%s' to complete ...", label)))(&resourceCleanup); err != nil {
		return err
	}

	ctx.executor(func() {
		defer resourceCleanup()
		fn()
	})

	return nil
}

func NewChannel[T any](ctx *Context, bufferCapacity uint) (chan T, func(), error) {

	var resourceCleanup func()

	if err := u.Get(NewResource(ctx, fmt.Sprintf("Channel: Waiting for channel (type %s, size: %d) to be closed ...", reflect.TypeFor[T]().Name(), bufferCapacity)))(&resourceCleanup); err != nil {
		return nil, nil, err
	}

	channel := make(chan T, bufferCapacity)

	channelCleanup := sync.OnceFunc(func() {
		close(channel)
		resourceCleanup()
	})

	return channel, channelCleanup, nil
}

func NewResource(ctx *Context, message string) (func(), error) {

	defer u.Lock(&ctx.lock)()

	if ctx.disposed {
		return nil, ErrDisposed
	}

	return ctx.resources.Add(ContextResource{message: message, stack: string(debug.Stack())}), nil
}
