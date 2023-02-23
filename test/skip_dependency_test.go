package di_test

import (
	"testing"

	di "coditory.com/goiku-di"
	"github.com/stretchr/testify/suite"
)

type SkipDependencySuite struct {
	suite.Suite
}

func (suite *SkipDependencySuite) TestReturnSkipError() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() (*Foo, error) { return nil, di.ErrSkipped })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Equal(di.ErrMissingObject, err)
}

func (suite *SkipDependencySuite) TestPanicWithSkipError() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo { panic(di.ErrSkipped) })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Equal(di.ErrMissingObject, err)
}

func (suite *SkipDependencySuite) TestResolveSliceWithNonSkippedDependnecies() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), &bar)
	ctxb.AddAs(new(Baz), func() *Foo { panic(di.ErrSkipped) })
	ctxb.Add(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(err)
	suite.NotNil(result)
	suite.NotNil(result.baz)
	suite.Equal(1, len(result.baz))
	suite.Equal(&bar, result.baz[0])
}

func TestSkipDependencySuite(t *testing.T) {
	suite.Run(t, new(SkipDependencySuite))
}
