package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type SkipDependencySuite struct {
	suite.Suite
}

func (suite *SkipDependencySuite) TestReturnSkipError() {
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() (*Foo, error) {
		return nil, di.ErrSkippedDependency
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(result)
	suite.Equal("missing dependency *di_test.Foo", err.Error())
}

func (suite *SkipDependencySuite) TestPanicWithSkipError() {
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo {
		panic(di.ErrSkippedDependency)
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(result)
	suite.Equal("missing dependency *di_test.Foo", err.Error())
}

func (suite *SkipDependencySuite) TestResolveSliceWithNonSkippedDependnecies() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.ProvideAs(new(Baz), func() *Foo {
		panic(di.ErrSkippedDependency)
	})
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Provide(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result := di.Get[*Boo](ctx)
	suite.NotNil(result)
	suite.NotNil(result.baz)
	suite.Equal(1, len(result.baz))
	suite.Equal(&bar, result.baz[0])
}

func TestSkipDependencySuite(t *testing.T) {
	suite.Run(t, new(SkipDependencySuite))
}
