package xhdl

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {

	res := ""

	err := Run(func(ctx Context) {
		res = "done"
	})

	assert.Equal(t, "done", res)
	assert.NoError(t, err)
}

func TestDirectErr(t *testing.T) {

	err := Run(func(ctx Context) {
		ctx.Throw(fmt.Errorf("error"))
	})

	assert.Error(t, err)
}

func TestInDirectErr(t *testing.T) {

	err := Run(func(ctx Context) {
		ind1(ctx)
	})

	assert.Error(t, err)
}

func ind1(ctx Context) {
	ind2(ctx)
}

func ind2(ctx Context) {
	ctx.Throw(fmt.Errorf("error"))
}

func TestNested(t *testing.T) {

	res := ""
	err := Run(func(ctx Context) {
		res = nested(ctx)
	})

	assert.NoError(t, err)
	assert.Equal(t, "error", res)

}

func nested(ctx Context) (res string) {
	err := Run(func(ctx Context) {
		ind1(ctx)
	})

	return err.Error()
}

func TestCallExternalWithoutIf(t *testing.T) {

	err := Run(func(ctx Context) {
		nerr := externalfn("")
		ctx.Throw(nerr)
	})

	assert.NoError(t, err)
}

func externalfn(errmsg string) error {
	if errmsg == "" {
		return nil
	} else {
		return fmt.Errorf(errmsg)
	}
}

var deferCalled int

func TestDeferedMethodsAreCalled(t *testing.T) {
	deferCalled = 0
	err := Run(func(ctx Context) {
		f1(ctx)
	})

	assert.Error(t, err)
	assert.Equal(t, 1, deferCalled)
}

func f1(ctx Context) {
	defer func() {
		deferCalled = 1
	}()
	f2(ctx)
}

func f2(ctx Context) {
	ctx.Throw(fmt.Errorf("Error!"))
}

func TestNoPanicsAreSwallowed(t *testing.T) {
	assert.Panics(t, func() {
		Run(func(ctx Context) {
			panic("ERROR")
		})
	})
}

func TestNoPanicsAreSwallowed2(t *testing.T) {
	assert.Panics(t, func() {
		Run(func(ctx Context) {
			var s []string
			_ = s[0:1]
		})
	})
}

func TestGetContextValid(t *testing.T) {
	err := Run(func(ctx Context) {

		// add cancel value to context
		ctx2, cancel := context.WithCancel(ctx)
		defer cancel()

		// and get the xhdl.Context from the WithCancel context
		xctx2 := GetContext(ctx2)
		assert.NotNil(t, xctx2)
	})

	assert.NoError(t, err)
}

func TestGetContextInValidShouldPanic(t *testing.T) {

	assert.Panics(t, func() {

		// this context is no valid xhdl context, so it should panic
		GetContext(context.Background())
	})

}
