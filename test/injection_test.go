package di_test

import (
	"testing"

	di "github.com/coditory/go-di"
	"github.com/stretchr/testify/suite"
)

type InjectionSuite struct {
	suite.Suite
}

func (suite *InjectionSuite) TestInjectParams() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Add(func(pfoo *Foo, pbar *Bar) *Boo { return &Boo{foo: pfoo, bar: pbar} })
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *InjectionSuite) TestInjectCastedParam() {
	type Boo struct {
		baz Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.Add(func(baz Baz) *Boo { return &Boo{baz: baz} })
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(&foo, result.baz)
}

func (suite *InjectionSuite) TestInjectContext() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Add(func(ctx *di.Context) *Boo {
		return &Boo{
			foo: di.GetOrPanic[*Foo](ctx),
			bar: di.GetOrPanic[*Bar](ctx),
		}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *InjectionSuite) TestInjectMixed() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Add(func(ctx *di.Context, foo *Foo) *Boo {
		return &Boo{
			foo: foo,
			bar: di.GetOrPanic[*Bar](ctx),
		}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *InjectionSuite) TestInjectMissingParam() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(func(foo *Foo, bar *Bar) *Boo {
		return &Boo{foo: foo, bar: bar}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(result)
	suite.NotNil(err)
	suite.Equal("missing object", err.Error())
}

func (suite *InjectionSuite) TestInjectSliceOfInterfaces() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Add(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.NotNil(result.baz)
	suite.Equal(2, len(result.baz))
	suite.Equal(&foo, result.baz[0])
	suite.Equal(&bar, result.baz[1])
}

func (suite *InjectionSuite) TestInjectSliceOfStructs() {
	type Boo struct {
		foo []Foo
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(Foo{name: "first"})
	ctxb.Add(Foo{name: "second"})
	ctxb.Add(func(foo []Foo) *Boo {
		return &Boo{foo: foo}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.NotNil(result.foo)
	suite.Equal(2, len(result.foo))
	suite.Equal("first", result.foo[0].name)
	suite.Equal("second", result.foo[1].name)
}

func (suite *InjectionSuite) TestInjectSliceOfStructPtrs() {
	type Boo struct {
		foo []*Foo
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&Foo{name: "first"})
	ctxb.Add(&Foo{name: "second"})
	ctxb.Add(func(foo []*Foo) *Boo {
		return &Boo{foo: foo}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.NotNil(result.foo)
	suite.Equal(2, len(result.foo))
	suite.Equal("first", result.foo[0].name)
	suite.Equal("second", result.foo[1].name)
}

func TestInjectionSuite(t *testing.T) {
	suite.Run(t, new(InjectionSuite))
}
