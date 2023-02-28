package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type ReadmeSamplesSuite struct {
	suite.Suite
}

func (suite *ReadmeSamplesSuite) TestDependencyRegistration() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(bar)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.GetOrPanic[*Foo](ctx))
	suite.Equal(bar, di.GetOrPanic[Bar](ctx))
}

type Named interface {
	Name() string
}

type NamedFoo struct {
	name string
}

func (f *NamedFoo) Name() string {
	return f.name
}

type NamedBar struct {
	name string
}

func (f *NamedBar) Name() string {
	return f.name
}

func (suite *ReadmeSamplesSuite) TestDependencyRetrievalByInterface() {
	foo := NamedFoo{name: "namedFoo"}
	foo2 := NamedFoo{name: "namedFoo2"}
	bar := NamedBar{name: "namedBar"}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Named), &foo)
	ctxb.AddAs(new(Named), &foo2)
	ctxb.Add(&foo) // can be registered multiple times
	ctxb.AddAs(new(Named), &bar)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.GetOrPanic[Named](ctx))
	suite.Equal([]Named{&foo, &foo2, &bar}, di.GetAll[Named](ctx))
}

func (suite *ReadmeSamplesSuite) TestNamedDependencyRegistration() {
	foo := Foo{name: "foo"}
	foo1 := Foo{name: "foo1"}
	foo2 := Foo{name: "foo2"}
	ctxb := di.NewContextBuilder()
	ctxb.AddNamed("foo1", &foo1)
	// ctxb.AddNamed("foo1", &Foo{name: "foo1-bis"}) // error
	ctxb.AddNamed("foo2", &foo2)
	ctxb.Add(&foo)
	ctx := ctxb.Build()
	suite.Equal(&foo1, di.GetOrPanic[*Foo](ctx))
	suite.Equal(&foo1, di.GetNamed[*Foo](ctx, "foo1"))
	suite.Equal(&foo2, di.GetNamed[*Foo](ctx, "foo2"))
	suite.Equal([]*Foo{&foo1, &foo2, &foo}, di.GetAll[*Foo](ctx))
}

func TestReadmeSamplesSuite(t *testing.T) {
	suite.Run(t, new(ReadmeSamplesSuite))
}
