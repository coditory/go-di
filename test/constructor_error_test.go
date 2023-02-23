package di_test

import (
	"testing"

	di "coditory.com/goiku-di"
	"github.com/stretchr/testify/suite"
)

type ConstructorErrorSuite struct {
	suite.Suite
}

func (suite *ConstructorErrorSuite) TestErrorResult() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() (*Foo, error) { return nil, errSimulated })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Equal(errSimulated, err)
}

func (suite *ConstructorErrorSuite) TestPanic() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo { panic(errSimulated) })
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Equal(errSimulated, err)
}

func (suite *SkipDependencySuite) TestPropagateErrorOnDependencySlice() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), func() *Foo { panic(errSimulated) })
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Add(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(result)
	suite.Equal(errSimulated, err)
}

func TestConstructorErrorSuite(t *testing.T) {
	suite.Run(t, new(ConstructorErrorSuite))
}
