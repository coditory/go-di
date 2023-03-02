package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type ParameterInjectionSuite struct {
	suite.Suite
}

func (suite *ParameterInjectionSuite) TestInjectParams() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Provide(func(pfoo *Foo, pbar *Bar) *Boo {
		return &Boo{foo: pfoo, bar: pbar}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *ParameterInjectionSuite) TestInjectCastedParam() {
	type Boo struct {
		baz Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.Provide(func(baz Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.Equal(&foo, result.baz)
}

func (suite *ParameterInjectionSuite) TestInjectContext() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Provide(func(ctx *di.Context) *Boo {
		return &Boo{
			foo: di.Get[*Foo](ctx),
			bar: di.Get[*Bar](ctx),
		}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *ParameterInjectionSuite) TestInjectMixed() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Add(&bar)
	ctxb.Provide(func(ctx *di.Context, foo *Foo) *Boo {
		return &Boo{
			foo: foo,
			bar: di.Get[*Bar](ctx),
		}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.Equal(&foo, result.foo)
	suite.Equal(&bar, result.bar)
}

func (suite *ParameterInjectionSuite) TestInjectMissingParam() {
	type Boo struct {
		foo *Foo
		bar *Bar
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo)
	ctxb.Provide(func(foo *Foo, bar *Bar) *Boo {
		return &Boo{foo: foo, bar: bar}
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Boo](ctx)
	suite.Nil(result)
	suite.NotNil(err)
	suite.Equal("could not create dependency *di_test.Boo, cause:\nmissing dependency *di_test.Bar", err.Error())
}

func (suite *ParameterInjectionSuite) TestInjectSliceOfInterfaces() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &foo)
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Provide(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.NotNil(result.baz)
	suite.Equal(2, len(result.baz))
	suite.Equal(&foo, result.baz[0])
	suite.Equal(&bar, result.baz[1])
}

func (suite *ParameterInjectionSuite) TestInjectSliceOfStructs() {
	type Boo struct {
		foo []Foo
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(Foo{id: "first"})
	ctxb.Add(Foo{id: "second"})
	ctxb.Provide(func(foo []Foo) *Boo {
		return &Boo{foo: foo}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.NotNil(result.foo)
	suite.Equal(2, len(result.foo))
	suite.Equal("first", result.foo[0].id)
	suite.Equal("second", result.foo[1].id)
}

func (suite *ParameterInjectionSuite) TestInjectSliceOfStructPtrs() {
	type Boo struct {
		foo []*Foo
	}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&Foo{id: "first"})
	ctxb.Add(&Foo{id: "second"})
	ctxb.Provide(func(foo []*Foo) *Boo {
		return &Boo{foo: foo}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.NotNil(result.foo)
	suite.Equal(2, len(result.foo))
	suite.Equal("first", result.foo[0].id)
	suite.Equal("second", result.foo[1].id)
}

func TestParameterInjectionSuite(t *testing.T) {
	suite.Run(t, new(ParameterInjectionSuite))
}
