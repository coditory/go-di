package di_test

import (
	"testing"

	di "coditory.com/goiku-di"
	"github.com/stretchr/testify/assert"
)

type cyclicFoo struct {
	bar *cyclicBar
}

type cyclicBar struct {
	baz *cyclicBaz
}

type cyclicBaz struct {
	foo *cyclicFoo
}

func provideCyclicFoo(bar *cyclicBar) *cyclicFoo {
	return &cyclicFoo{bar: bar}
}

func provideCyclicBar(baz *cyclicBaz) *cyclicBar {
	return &cyclicBar{baz: baz}
}

func provideCyclicBaz(foo *cyclicFoo) *cyclicBaz {
	return &cyclicBaz{foo: foo}
}

func Test_CyclicDependency(t *testing.T) {
	ctxb := di.NewContextBuilder()
	ctxb.Add(provideCyclicFoo)
	ctxb.Add(provideCyclicBar)
	ctxb.Add(provideCyclicBaz)
	ctx := ctxb.Build()
	result, err := di.Get[*cyclicFoo](ctx)
	assert.ErrorIs(t, err, di.ErrDependencyCycle)
	assert.Nil(t, result)
}
