package xhdl

import "context"

type Context interface {
	context.Context
	Throw(err error)
	RunNested(f func(xctx Context)) error
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

func (xc *xcontext) RunNested(f func(xctx Context)) error {
	return RunContext(xc.Context, f)
}

func NewContext(f func(xctx Context)) func(ctx context.Context) error {

	return func(ctx context.Context) (err error) {
		xctx := &xcontext{ctx}

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

func Run(f func(xctx Context)) (err error) {
	return RunContext(context.Background(), f)
}

func RunContext(ctx context.Context, f func(xctx Context)) (err error) {
	instance := NewContext(f)
	return instance(ctx)
}
