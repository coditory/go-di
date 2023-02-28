package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type NamedDependencySuite struct {
	suite.Suite
}

func (suite *NamedDependencySuite) TestGetNamedDependencyByName() {
	foo1 := Foo{name: "foo1"}
	foo2 := Foo{name: "foo2"}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&Foo{name: "foo"})
	ctxb.AddNamed("foo1", &foo1)
	ctxb.AddNamed("foo2", &foo2)
	ctx := ctxb.Build()
	result, err := di.GetNamedOrErr[*Foo](ctx, "foo1")
	suite.Nil(err)
	suite.Equal(&foo1, result)
	result, err = di.GetNamedOrErr[*Foo](ctx, "foo2")
	suite.Nil(err)
	suite.Equal(&foo2, result)
}

func (suite *NamedDependencySuite) TestGetNamedDependencyByType() {
	foo1 := Foo{name: "foo1"}
	foo2 := Foo{name: "foo2"}
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo1", &foo1)
	ctxb.AddNamed("foo2", &foo2)
	ctxb.Add(&foo)
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(err)
	suite.Equal(&foo1, result)
	all, err := di.GetAllOrErr[*Foo](ctx)
	suite.Nil(err)
	suite.Equal([]*Foo{&foo1, &foo2, &foo}, all)
}

func (suite *NamedDependencySuite) TestErrorOnDuplicatedName() {
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo", &Foo{name: "foo1"})
	err := ctxb.AddNamedOrErr("foo", &Foo{name: "foo2"})
	suite.Equal("duplicated dependency name: foo", err.Error())
}

func (suite *NamedDependencySuite) TestErrorOnInvalidType() {
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo", &foo)
	ctx := ctxb.Build()
	obj, err := di.GetNamedOrErr[*Bar](ctx, "foo")
	suite.Nil(obj)
	suite.Equal("could not cast *di_test.Foo (name: foo) to *di_test.Bar", err.Error())
}

func TestNamedDependencySuite(t *testing.T) {
	suite.Run(t, new(NamedDependencySuite))
}
