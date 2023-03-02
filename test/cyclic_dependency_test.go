package di_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
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

func provideCyclicFooWithCtx(ctx *di.Context) *cyclicFoo {
	return &cyclicFoo{bar: di.GetOrPanic[*cyclicBar](ctx)}
}

func provideCyclicBarWithCtx(ctx *di.Context) *cyclicBar {
	return &cyclicBar{baz: di.GetOrPanic[*cyclicBaz](ctx)}
}

func provideCyclicBazWithCtx(ctx *di.Context) *cyclicBaz {
	return &cyclicBaz{foo: di.GetOrPanic[*cyclicFoo](ctx)}
}

type CyclicDependencySuite struct {
	suite.Suite
}

func (suite *CyclicDependencySuite) TestCyclicDependencyWithInjection() {
	tests := []struct {
		title    string
		register func(*di.ContextBuilder)
	}{
		{
			title: "param",
			register: func(ctxb *di.ContextBuilder) {
				ctxb.Add(provideCyclicFoo)
				ctxb.Add(provideCyclicBar)
				ctxb.Add(provideCyclicBaz)
			},
		},
		{
			title: "context",
			register: func(ctxb *di.ContextBuilder) {
				ctxb.Add(provideCyclicFooWithCtx)
				ctxb.Add(provideCyclicBarWithCtx)
				ctxb.Add(provideCyclicBazWithCtx)
			},
		},
		{
			title: "mixed",
			register: func(ctxb *di.ContextBuilder) {
				ctxb.Add(provideCyclicFoo)
				ctxb.Add(provideCyclicBarWithCtx)
				ctxb.Add(provideCyclicBaz)
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.title, func() {
			ctxb := di.NewContextBuilder()
			tt.register(ctxb)
			ctx := ctxb.Build()
			result, err := di.GetOrErr[*cyclicFoo](ctx)
			suite.Nil(result)
			suite.Equal(strings.Join([]string{
				"could not create dependency *di_test.cyclicFoo, cause:",
				"could not create dependency *di_test.cyclicBar, cause:",
				"could not create dependency *di_test.cyclicBaz, cause:",
				"cyclic dependency: *di_test.cyclicFoo -> *di_test.cyclicBar -> *di_test.cyclicBaz -> *di_test.cyclicFoo",
			}, "\n"),
				err.Error())
			suite.Equal(di.ErrTypeDependencyCreation, err.ErrType())
			suite.Equal(di.ErrTypeCyclicDependency, err.RootCause().(*di.Error).ErrType())
		})
	}
}

func (suite *CyclicDependencySuite) TestCyclicDependencyOfOne() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func(foo *Foo) *Foo {
		return foo
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(result)
	suite.Equal("could not create dependency *di_test.Foo, cause:\ncyclic dependency: *di_test.Foo -> *di_test.Foo", err.Error())
	suite.Equal(di.ErrTypeDependencyCreation, err.ErrType())
	suite.Equal(di.ErrTypeCyclicDependency, err.RootCause().(*di.Error).ErrType())
}

func (suite *CyclicDependencySuite) TestNoErrorThrownWhenNotRetrieved() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func(foo *Foo) *Foo {
		return foo
	})
	ctxb.Add(func() *Bar {
		return &bar
	})
	ctx := ctxb.Build()
	result := di.Get[*Bar](ctx)
	suite.Equal(&bar, result)
}

func TestCyclicDependencySuite(t *testing.T) {
	suite.Run(t, new(CyclicDependencySuite))
}
