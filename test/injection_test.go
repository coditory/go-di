package di_test

import (
	"testing"

	di "coditory.com/goiku-di"
	"github.com/stretchr/testify/suite"
)

type InjectionSuite struct {
	suite.Suite
}

func (suite *InjectionSuite) TestInjectWithParams() {
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

func TestInjectionSuite(t *testing.T) {
	suite.Run(t, new(InjectionSuite))
}
