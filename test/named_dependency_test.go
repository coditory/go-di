package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type NamedDependencySuite struct {
	suite.Suite
}

func (suite *LifecycleSuite) TestGetNamedDependencyByName() {
	foo1 := Foo{id: "foo1"}
	foo2 := Foo{id: "foo2"}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&Foo{id: "foo"})
	ctxb.AddNamed("foo1", &foo1)
	ctxb.AddNamed("foo2", &foo2)
	ctx := ctxb.Build()
	result := di.GetNamed[*Foo](ctx, "foo1")
	suite.Equal(&foo1, result)
	result = di.GetNamed[*Foo](ctx, "foo2")
	suite.Equal(&foo2, result)
}

func (suite *LifecycleSuite) TestGetNamedDependencyByType() {
	foo1 := Foo{id: "foo1"}
	foo2 := Foo{id: "foo2"}
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo1", &foo1)
	ctxb.AddNamed("foo2", &foo2)
	ctxb.Add(&foo)
	ctx := ctxb.Build()
	result := di.Get[*Foo](ctx)
	suite.Equal(&foo1, result)
	all := di.GetAll[*Foo](ctx)
	suite.Equal([]*Foo{&foo1, &foo2, &foo}, all)
}

func (suite *LifecycleSuite) TestRegisterNamedDependencyTwiceForDifferentTypes() {
	foo1 := Foo{id: "foo1"}
	ctxb := di.NewContextBuilder()
	ctxb.AddNamedAs("foo", new(any), &foo1)
	ctxb.AddNamed("foo", &foo1)
	ctx := ctxb.Build()
	result := di.Get[*Foo](ctx)
	suite.Equal(&foo1, result)
	resultAny := di.Get[any](ctx)
	suite.Equal(&foo1, resultAny)
	all := di.GetAll[*Foo](ctx)
	suite.Equal([]*Foo{&foo1}, all)
}

func (suite *LifecycleSuite) TestErrorOnDuplicatedName() {
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo", &Foo{id: "foo1"})
	err := ctxb.AddNamedOrErr("foo", &Foo{id: "foo2"})
	suite.Equal("duplicated dependency name: foo", err.Error())
}

func (suite *LifecycleSuite) TestErrorOnInvalidType() {
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo", &foo)
	ctx := ctxb.Build()
	obj, err := di.GetNamedOrErr[*Bar](ctx, "foo")
	suite.Nil(obj)
	suite.Equal("could not cast *di_test.Foo (name: foo) to *di_test.Bar", err.Error())
}

func TestNamedDependencySuite(t *testing.T) {
	suite.Run(t, new(LifecycleSuite))
}
