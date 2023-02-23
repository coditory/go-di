package di_test

import (
	"testing"

	di "github.com/coditory/go-di"
	"github.com/stretchr/testify/suite"
)

type SingleInitSuite struct {
	suite.Suite
}

func (suite *SingleInitSuite) TestMultipleGet() {
	inits := 0
	ctxb := di.NewContextBuilder()
	ctxb.Add(func() *Foo { inits++; return &foo })
	ctx := ctxb.Build()
	_, _ = di.Get[*Foo](ctx)
	_, _ = di.Get[*Foo](ctx)
	_, _ = di.GetAll[*Foo](ctx)
	suite.Equal(1, inits)
}

func (suite *SingleInitSuite) TestMultipleAddDifferentTypes() {
	inits := 0
	ctor := func() *Foo { inits++; return &Foo{} }
	ctxb := di.NewContextBuilder()
	ctxb.Add(ctor)
	ctxb.AddAs(new(Baz), ctor)
	ctx := ctxb.Build()
	rfoo, _ := di.Get[*Foo](ctx)
	rbaz, _ := di.Get[Baz](ctx)
	suite.Equal(1, inits)
	suite.Equal(rbaz, rfoo)
}

func TestSingleInitSuite(t *testing.T) {
	suite.Run(t, new(SingleInitSuite))
}
