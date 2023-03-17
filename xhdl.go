package xhdl

import (
	"context"
	"fmt"
)

// Context is the interface that allows err to be thrown. You have to pass this
// to all func participating in the xhdl Context.
type Context interface {
	context.Context

	// Throw causes the xhdl Context to abort, and to return the err to the caller.
	// * Throw check for err != nil internally, and is a nop if err == nil
	// * Even if the go compiler knows that err != nil, there may be a 'missing return' error, because it can't look into the interface func implementation
	Throw(err error)

	// RunNested allows a inner function to perform it's own error handling.
	// This is just a shortcut to RunContext
	RunNested(f func(ctx Context)) error
}

type xcontext struct {
	context.Context
}

type wrappedErr struct {
	err error
}

func (xc *xcontext) Throw(err error) {
	if err != nil {
		panic(&wrappedErr{err})
	}
}

func (xc *xcontext) RunNested(f func(ctx Context)) error {
	return RunContext(xc.Context, f)
}

type ctxkeytype int

const ctxkey = ctxkeytype(0)

// GetContext returns a xhdl.Context interface for the given context.Context
// This will panic if ctx or a parent context if not a xhdl.Context.
// Use this func you want to "get back" the xhdl Context that is hidden behind
// a Context returned by functions like context.WithCancel, .WithDeadline or .WithValue
func GetContext(ctx context.Context) Context {

	// validate context by value
	xctx, valid := ctx.Value(ctxkey).(*xcontext)
	if !valid {
		panic(fmt.Errorf("the provivided context is not managed by xhdl"))
	}

	return xctx
}

// NewContext creates returns a func that can later be executed with xhdl
func NewContext(f func(ctx Context)) func(ctx context.Context) error {

	return func(ctx context.Context) (err error) {
		xctx := &xcontext{}
		xctx.Context = context.WithValue(ctx, ctxkey, xctx)

		defer func() {
			if r := recover(); r != nil {
				lerr, iserr := r.(*wrappedErr)
				if !iserr {
					panic(r)
				} else {
					err = lerr.err
				}
			}
		}()

		f(xctx)
		return
	}
}

// Run executes f in a xhdl Context and returns err if ctx.Throw() has been
// call with a non-nil err.
func Run(f func(ctx Context)) (err error) {
	return RunContext(context.Background(), f)
}

// RunContext executes f in a xhdl Context based on a parent context and returns err if ctx.Throw() has been
// call with a non-nil err.
func RunContext(ctx context.Context, f func(ctx Context)) (err error) {
	instance := NewContext(f)
	return instance(ctx)
}
