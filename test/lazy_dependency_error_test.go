package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type LazyDependencyErrorSuite struct {
	suite.Suite
}

func (suite *LazyDependencyErrorSuite) TestErrorResult() {
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() (*Foo, error) {
		return nil, errSimulated
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(result)
	suite.Equal("could not create dependency *di_test.Foo, cause:\nsimulated", err.Error())
	suite.Equal(di.ErrTypeDependencyCreation, err.ErrType())
	suite.Equal(errSimulated, err.RootCause())
	suite.ErrorIs(err, errSimulated)
}

func (suite *LazyDependencyErrorSuite) TestPanic() {
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo {
		panic(errSimulated)
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Foo](ctx)
	suite.Nil(result)
	suite.Equal(di.ErrTypeDependencyCreation, err.ErrType())
	suite.Equal(errSimulated, err.RootCause())
	suite.ErrorIs(err, errSimulated)
}

func (suite *LazyDependencyErrorSuite) TestErrorOnSliceDependency() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.ProvideAs(new(Baz), func() *Foo {
		panic(errSimulated)
	})
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Provide(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*Boo](ctx)
	suite.Nil(result)
	suite.Equal("could not create dependency *di_test.Boo, cause:\ncould not create dependency di_test.Baz, cause:\nsimulated", err.Error())
	suite.Equal(di.ErrTypeDependencyCreation, err.ErrType())
	suite.ErrorIs(err, errSimulated)
}

func TestLazyDependencyErrorSuite(t *testing.T) {
	suite.Run(t, new(LazyDependencyErrorSuite))
}
