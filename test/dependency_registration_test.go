package di_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type DependencyRegistrationSuite struct {
	suite.Suite
}

func (suite *DependencyRegistrationSuite) TestGetByType() {
	tests := []struct {
		value any
		get   func(ctx *di.Context) (any, error)
	}{
		{
			value: &foo,
			get:   func(ctx *di.Context) (any, error) { return di.GetOrErr[*Foo](ctx) },
		},
		{
			value: foo,
			get:   func(ctx *di.Context) (any, error) { return di.GetOrErr[Foo](ctx) },
		},
		{
			value: 42,
			get:   func(ctx *di.Context) (any, error) { return di.GetOrErr[int](ctx) },
		},
		{
			value: "text",
			get:   func(ctx *di.Context) (any, error) { return di.GetOrErr[string](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.Add(tt.value)
			ctx := ctxb.Build()
			result, err := tt.get(ctx)
			suite.Nil(err)
			suite.Equal(tt.value, result)
		})
	}
}

func (suite *DependencyRegistrationSuite) TestGetByInterface() {
	tests := []struct {
		value any
		iface any
		get   func(ctx *di.Context) any
	}{
		{
			value: &foo,
			iface: new(Baz),
			get:   func(ctx *di.Context) any { return di.Get[Baz](ctx) },
		},
		{
			value: (*Foo)(nil),
			iface: new(Baz),
			get:   func(ctx *di.Context) any { return di.Get[Baz](ctx) },
		},
		{
			value: bar,
			iface: new(Baz),
			get:   func(ctx *di.Context) any { return di.Get[Baz](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.AddAs(tt.iface, tt.value)
			ctx := ctxb.Build()
			result := tt.get(ctx)
			suite.Equal(tt.value, result)
		})
	}
}

func (suite *DependencyRegistrationSuite) TestGetAllByType() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(foo1)
	ctxb.Add(foo2)
	ctx := ctxb.Build()
	result := di.GetAll[*Foo](ctx)
	suite.Equal([]*Foo{foo1, foo2}, result)
}

func (suite *DependencyRegistrationSuite) TestGetAllByInterface() {
	foo1 := &Foo{}
	foo2 := &Foo{}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), foo1)
	ctxb.AddAs(new(Baz), foo2)
	ctx := ctxb.Build()
	result := di.GetAll[Baz](ctx)
	suite.Equal([]Baz{foo1, foo2}, result)
}

func (suite *DependencyRegistrationSuite) TestGetAllMissing() {
	ctxb := di.NewContextBuilder()
	ctx := ctxb.Build()
	result, err := di.GetOrErr[Baz](ctx)
	suite.Nil(result)
	suite.Equal("missing dependency di_test.Baz", err.Error())
	suite.IsType(new(di.MissingDependencyError), err)
}

func (suite *DependencyRegistrationSuite) TestGetMissing() {
	ctxb := di.NewContextBuilder()
	ctx := ctxb.Build()
	result := di.GetAll[Baz](ctx)
	suite.Equal([]Baz{}, result)
}

func TestDependencyRegistrationSuite(t *testing.T) {
	suite.Run(t, new(DependencyRegistrationSuite))
}
