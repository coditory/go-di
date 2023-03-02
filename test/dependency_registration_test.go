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
		get   func(ctx *di.Context) any
	}{
		{
			value: &foo,
			get:   func(ctx *di.Context) any { return di.Get[*Foo](ctx) },
		},
		{
			value: foo,
			get:   func(ctx *di.Context) any { return di.Get[Foo](ctx) },
		},
		{
			value: 42,
			get:   func(ctx *di.Context) any { return di.Get[int](ctx) },
		},
		{
			value: "text",
			get:   func(ctx *di.Context) any { return di.Get[string](ctx) },
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.Add(tt.value)
			ctx := ctxb.Build()
			result := tt.get(ctx)
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

func someFunc() string {
	return "result from someFunc()"
}

func someFunc2() string {
	return "result from someFunc2()"
}

func someFunc3() int {
	return 42
}

func someFunc4(a int) string {
	return fmt.Sprintf("result from someFunc4(%d)", a)
}

func (suite *DependencyRegistrationSuite) TestGetByFunc() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(someFunc)
	ctxb.Add(someFunc3)
	ctxb.Add(someFunc4)
	ctx := ctxb.Build()
	result := di.Get[func() string](ctx)
	suite.Equal(someFunc(), result())
	result2 := di.Get[func() int](ctx)
	suite.Equal(someFunc3(), result2())
	result3 := di.Get[func(int) string](ctx)
	suite.Equal(someFunc4(11), result3(11))
}

func (suite *DependencyRegistrationSuite) TestGetAllByFunc() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(someFunc)
	ctxb.Add(someFunc2)
	ctx := ctxb.Build()
	result := di.GetAll[func() string](ctx)
	suite.Equal(2, len(result))
	suite.Equal(someFunc(), result[0]())
	suite.Equal(someFunc2(), result[1]())
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
	suite.Equal(di.ErrTypeMissingDependency, err.ErrType())
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
