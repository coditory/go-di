package di_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type LazyDependencySuite struct {
	suite.Suite
}

func (suite *LazyDependencySuite) TestGetByType() {
	tests := []struct {
		value   any
		provide func(ctxb *di.ContextBuilder)
		get     func(ctx *di.Context) (any, error)
	}{
		{
			value:   &foo,
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() *Foo { return &foo }) },
			get:     func(ctx *di.Context) (any, error) { return di.GetOrErr[*Foo](ctx) },
		},
		{
			value:   foo,
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() Foo { return foo }) },
			get:     func(ctx *di.Context) (any, error) { return di.GetOrErr[Foo](ctx) },
		},
		{
			value:   42,
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() int { return 42 }) },
			get:     func(ctx *di.Context) (any, error) { return di.GetOrErr[int](ctx) },
		},
		{
			value:   "text",
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() string { return "text" }) },
			get:     func(ctx *di.Context) (any, error) { return di.GetOrErr[string](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s %+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			tt.provide(ctxb)
			ctx := ctxb.Build()
			result, err := tt.get(ctx)
			suite.Nil(err)
			suite.Equal(tt.value, result)
		})
	}
}

func (suite *LazyDependencySuite) TestGetByInterface() {
	tests := []struct {
		value   any
		iface   any
		provide func(ctxb *di.ContextBuilder)
		get     func(ctx *di.Context) any
	}{
		{
			value:   &foo,
			iface:   new(Baz),
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() Baz { return &foo }) },
			get:     func(ctx *di.Context) any { return di.Get[Baz](ctx) },
		},
		{
			value:   bar,
			iface:   new(Baz),
			provide: func(ctxb *di.ContextBuilder) { ctxb.Provide(func() Baz { return bar }) },
			get:     func(ctx *di.Context) any { return di.Get[Baz](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			tt.provide(ctxb)
			ctx := ctxb.Build()
			result := tt.get(ctx)
			suite.Equal(tt.value, result)
		})
	}
}

func (suite *LazyDependencySuite) TestGetAllByType() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo { return foo1 })
	ctxb.Provide(func() *Foo { return foo2 })
	ctx := ctxb.Build()
	result := di.GetAll[*Foo](ctx)
	suite.Equal([]*Foo{foo1, foo2}, result)
}

func (suite *LazyDependencySuite) TestGetAllByInterface() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.ProvideAs(new(Baz), func() *Foo { return foo1 })
	ctxb.ProvideAs(new(Baz), func() *Foo { return foo2 })
	ctx := ctxb.Build()
	result := di.GetAll[Baz](ctx)
	suite.Equal([]Baz{foo1, foo2}, result)
}

func TestLazyDependencySuite(t *testing.T) {
	suite.Run(t, new(LazyDependencySuite))
}
