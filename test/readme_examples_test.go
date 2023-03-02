package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/coditory/go-di"
)

type ReadmeExamplesSuite struct {
	suite.Suite
}

func (suite *ReadmeExamplesSuite) TestReadme01_FirstExample() {
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.AddAs(new(Baz), &foo2)
	ctxb.AddAs(new(Baz), &bar)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.Get[Baz](ctx))
	suite.Equal([]Baz{&foo, &foo2, &bar}, di.GetAll[Baz](ctx))
}

func (suite *ReadmeExamplesSuite) TestReadme02_ObjRefByType() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.Get[*Foo](ctx))
	suite.Equal([]*Foo{&foo, &foo2}, di.GetAll[*Foo](ctx))
}

func (suite *ReadmeExamplesSuite) TestReadme03_ObjByType() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(foo)
	ctxb.Add(foo2)
	ctx := ctxb.Build()
	suite.Equal(foo, di.Get[Foo](ctx))
	foos := di.GetAll[Foo](ctx)
	suite.Equal(2, len(foos))
	suite.Equal("foo", foos[0].Id())
	suite.Equal("foo2", foos[1].Id())
}

func (suite *ReadmeExamplesSuite) TestReadme04_ObjRefByIface() {
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.AddAs(new(Baz), &foo2)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.Get[Baz](ctx))
	suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
}

func (suite *ReadmeExamplesSuite) TestReadme05_ObjRefByIfaceAndType() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.AddAs(new(Baz), &foo)
	ctxb.Add(&foo2)
	ctxb.AddAs(new(Baz), &foo2)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.Get[Baz](ctx))
	suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
	suite.Equal(&foo, di.Get[*Foo](ctx))
	suite.Equal([]*Foo{&foo, &foo2}, di.GetAll[*Foo](ctx))
}

func (suite *ReadmeExamplesSuite) TestReadme06_ObjRefByName() {
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.AddNamedAs("special-foo", new(Baz), &foo2)
	ctx := ctxb.Build()
	suite.Equal(&foo, di.Get[Baz](ctx))
	suite.Equal(&foo2, di.GetNamed[Baz](ctx, "special-foo"))
	suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
}

func (suite *ReadmeExamplesSuite) TestReadme07_LazyDependencies() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo {
		return &foo
	})
	ctxb.Provide(func() *Bar {
		return &bar
	})
	ctxb.Provide(func(foo *Foo, bar *Bar) *Boo {
		return &Boo{foo: foo, bar: bar}
	})
	ctx := ctxb.Build()
	boo := di.Get[*Boo](ctx)
	suite.Equal(&foo, boo.foo)
	suite.Equal(&bar, boo.bar)
}

func (suite *ReadmeExamplesSuite) TestReadme08_LazyDependenciesWithCtx() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo {
		return &foo
	})
	ctxb.Provide(func() *Bar {
		return &bar
	})
	ctxb.Provide(func(ctx *di.Context) *Boo {
		return &Boo{foo: di.Get[*Foo](ctx), bar: di.Get[*Bar](ctx)}
	})
	ctx := ctxb.Build()
	boo := di.Get[*Boo](ctx)
	suite.Equal(&foo, boo.foo)
	suite.Equal(&bar, boo.bar)
}

func (suite *ReadmeExamplesSuite) TestReadme09_LazyDependenciesWithSingleCreation() {
	creations := 0
	createFoo := func() *Foo {
		creations++
		return &foo
	}
	ctxb := di.NewContextBuilder()
	ctxb.Provide(createFoo)
	ctxb.ProvideAs(new(Baz), createFoo)
	ctx := ctxb.Build()
	di.Get[Baz](ctx)
	di.Get[*Foo](ctx)
	suite.Equal(1, creations)
}

func TestReadmeExamplesSuite(t *testing.T) {
	suite.Run(t, new(ReadmeExamplesSuite))
}
