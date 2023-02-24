package di_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type EagerDependencySuite struct {
	suite.Suite
}

func (suite *EagerDependencySuite) TestGetByType() {
	tests := []struct {
		value any
		get   func(ctx *di.Context) (any, error)
	}{
		{
			value: &foo,
			get:   func(ctx *di.Context) (any, error) { return di.Get[*Foo](ctx) },
		},
		{
			value: foo,
			get:   func(ctx *di.Context) (any, error) { return di.Get[Foo](ctx) },
		},
		{
			value: 42,
			get:   func(ctx *di.Context) (any, error) { return di.Get[int](ctx) },
		},
		{
			value: "text",
			get:   func(ctx *di.Context) (any, error) { return di.Get[string](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.Add(tt.value)
			ctx := ctxb.Build()
			result, err := tt.get(ctx)
			suite.Nil(err, "received error")
			suite.Equal(tt.value, result, "retrieved value did not match")
		})
	}
}

func (suite *EagerDependencySuite) TestGetByInterface() {
	tests := []struct {
		value any
		iface any
		get   func(ctx *di.Context) (any, error)
	}{
		{
			value: &foo,
			iface: new(Baz),
			get:   func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
		},
		{
			value: (*Foo)(nil),
			iface: new(Baz),
			get:   func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
		},
		{
			value: bar,
			iface: new(Baz),
			get:   func(ctx *di.Context) (any, error) { return di.Get[Baz](ctx) },
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

func (suite *EagerDependencySuite) TestGetAllByType() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(foo1)
	ctxb.Add(foo2)
	ctx := ctxb.Build()
	result, err := di.GetAll[*Foo](ctx)
	suite.Nil(err)
	suite.Equal([]*Foo{foo1, foo2}, result)
}

func (suite *EagerDependencySuite) TestGetAllByInterface() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), foo1)
	ctxb.AddAs(new(Baz), foo2)
	ctx := ctxb.Build()
	result, err := di.GetAll[Baz](ctx)
	suite.Nil(err)
	suite.Equal([]Baz{foo1, foo2}, result)
}

func (suite *EagerDependencySuite) TestGetAllMissing() {
	ctxb := di.NewContextBuilder()
	ctx := ctxb.Build()
	result, err := di.Get[Baz](ctx)
	suite.Nil(result)
	suite.Equal(di.ErrMissingObject, err)
}

func (suite *EagerDependencySuite) TestGetMissing() {
	ctxb := di.NewContextBuilder()
	ctx := ctxb.Build()
	result, err := di.GetAll[Baz](ctx)
	suite.Nil(err)
	suite.Equal([]Baz{}, result)
}

func TestEagerDependencySuite(t *testing.T) {
	suite.Run(t, new(EagerDependencySuite))
}
