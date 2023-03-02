package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type SingleCreationSuite struct {
	suite.Suite
}

func (suite *SingleCreationSuite) TestMultipleGet() {
	inits := 0
	ctxb := di.NewContextBuilder()
	ctxb.Provide(func() *Foo { inits++; return &foo })
	ctx := ctxb.Build()
	_, _ = di.GetOrErr[*Foo](ctx)
	_, _ = di.GetOrErr[*Foo](ctx)
	_, _ = di.GetAllOrErr[*Foo](ctx)
	suite.Equal(1, inits)
}

func (suite *SingleCreationSuite) TestMultipleAddDifferentTypes() {
	inits := 0
	ctor := func() *Foo { inits++; return &Foo{} }
	ctxb := di.NewContextBuilder()
	ctxb.Provide(ctor)
	ctxb.ProvideAs(new(Baz), ctor)
	ctx := ctxb.Build()
	rfoo, _ := di.GetOrErr[*Foo](ctx)
	rbaz, _ := di.GetOrErr[Baz](ctx)
	suite.Equal(1, inits)
	suite.Equal(rbaz, rfoo)
}

func TestSingleCreationSuite(t *testing.T) {
	suite.Run(t, new(SingleCreationSuite))
}
