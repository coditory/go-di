package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type CreationErrorSuite struct {
	suite.Suite
}

func (suite *CreationErrorSuite) TestErrorResult() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() (*Foo, error) {
		return nil, errSimulated
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.Equal("could not create dependency *di_test.Foo, cause:\nsimulated", err.Error())
	suite.IsType(new(di.DependencyCreationError), err)
	suite.ErrorIs(err, errSimulated)
}

func (suite *CreationErrorSuite) TestPanic() {
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo {
		panic(errSimulated)
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Foo](ctx)
	suite.Nil(result)
	suite.IsType(new(di.DependencyCreationError), err)
	suite.ErrorIs(err, errSimulated)
}

func (suite *CreationErrorSuite) TestPropagateErrorOnDependencySlice() {
	type Boo struct {
		baz []Baz
	}
	ctxb := di.NewContextBuilder()
	ctxb.AddAs(new(Baz), func() *Foo {
		panic(errSimulated)
	})
	ctxb.AddAs(new(Baz), &bar)
	ctxb.Add(func(baz []Baz) *Boo {
		return &Boo{baz: baz}
	})
	ctx := ctxb.Build()
	result, err := di.Get[*Boo](ctx)
	suite.Nil(result)
	suite.Equal("could not create dependency *di_test.Boo, cause:\ncould not create dependency di_test.Baz, cause:\nsimulated", err.Error())
	suite.IsType(new(di.DependencyCreationError), err)
	suite.ErrorIs(err, errSimulated)
}

func TestCreationErrorSuite(t *testing.T) {
	suite.Run(t, new(CreationErrorSuite))
}
