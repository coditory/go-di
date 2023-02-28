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
	ctxb.Add(func() (*Foo, error) { return nil, di.ErrSkippedDependency })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Error(err, "abc")
}

func (suite *SkipDependencySuite) TestPanicWithSkipError() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo { panic(di.ErrSkippedDependency) })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Error(err, "abc")
}

func (suite *SkipDependencySuite) TestResolveSliceWithNonSkippedDependnecies() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), func() *Foo { panic(di.ErrSkippedDependency) })
	ctxb.AddAs(new(Baz), &bar)
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
