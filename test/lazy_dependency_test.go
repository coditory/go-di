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
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() *Foo { return &foo }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[*Foo](ctx) },
		},
		{
			value:   foo,
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() Foo { return foo }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[Foo](ctx) },
		},
		{
			value:   42,
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() int { return 42 }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[int](ctx) },
		},
		{
			value:   "text",
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() string { return "text" }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[string](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			tt.provide(ctxb)
			ctx := ctxb.Build()
			result, err := tt.get(ctx)
			suite.Nil(err, "received error")
			suite.Equal(tt.value, result, "retrieved value did not match")
		})
	}
}

func (suite *LazyDependencySuite) TestGetByInterface() {
	tests := []struct {
		value   any
		iface   any
		provide func(ctxb *di.ContextBuilder)
		get     func(ctx *di.Context) (any, error)
	}{
		{
			value:   &foo,
			iface:   new(Baz),
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() *Foo { return &foo }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
		},
		{
			value:   (*Foo)(nil),
			iface:   new(Baz),
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() *Foo { return nil }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
		},
		{
			value:   bar,
			iface:   new(Baz),
			provide: func(ctxb *di.ContextBuilder) { ctxb.Add(func() Bar { return bar }) },
			get:     func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.AddAs(tt.iface, tt.value)
			ctx := ctxb.Build()
			result, err := tt.get(ctx)
			suite.Nil(err, "received error")
			suite.Equal(tt.value, result, "retrieved value did not match")
		})
	}
}

func (suite *LazyDependencySuite) TestGetAllByType() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo { return foo1 })
	ctxb.Add(func() *Foo { return foo2 })
	ctx := ctxb.Build()
	result, err := di.GetAll[*Foo](ctx)
	suite.Nil(err)
	suite.Equal([]*Foo{foo1, foo2}, result)
}

func (suite *LazyDependencySuite) TestGetAllByInterface() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), func() *Foo { return foo1 })
	ctxb.AddAs(new(Baz), func() *Foo { return foo2 })
	ctx := ctxb.Build()
	result, err := di.GetAll[Baz](ctx)
	suite.Nil(err)
	suite.Equal([]Baz{foo1, foo2}, result)
}

func TestLazyDependencySuite(t *testing.T) {
	suite.Run(t, new(LazyDependencySuite))
}
