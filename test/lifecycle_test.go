package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	stdcontext "context"

	di "github.com/coditory/go-di"
)

type CtxAwareFoo struct {
	initialized     int
	shutdown        int
	errOnInitialize bool
	errOnShutdown   bool
}

func (f *CtxAwareFoo) Initialize() {
	if f.errOnInitialize {
		panic(errSimulated)
	}
	f.initialized++
}

func (f *CtxAwareFoo) Shutdown(context stdcontext.Context) {
	if f.errOnShutdown {
		panic(errSimulated)
	}
	f.shutdown++
}

type LifecycleSuite struct {
	suite.Suite
}

func (suite *LifecycleSuite) TestDependencyInit() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	suite.Equal(foo1.initialized, 0)
	suite.Equal(foo2.initialized, 0)
	ctx.Initialize()
	suite.Equal(foo1.initialized, 1)
	suite.Equal(foo2.initialized, 1)
}

func (suite *LifecycleSuite) TestDuplicatedInitError() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	ctx.Initialize()
	err := ctx.InitializeOrErr()
	suite.Equal("context lifecycle error: context already initialized", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeLifecycle)
}

func (suite *LifecycleSuite) TestDependencyInitError() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{errOnInitialize: true}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	err := ctx.InitializeOrErr()
	suite.Equal("could not initialize dependency: *di_test.CtxAwareFoo, cause:\nsimulated", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeDependencyInitialization)
}

func (suite *LifecycleSuite) TestDependencyShutdown() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	suite.Equal(foo1.shutdown, 0)
	suite.Equal(foo2.shutdown, 0)
	ctx.Shutdown(stdcontext.TODO())
	suite.Equal(foo1.shutdown, 1)
	suite.Equal(foo2.shutdown, 1)
}

func (suite *LifecycleSuite) TestDuplicatedShutdownError() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	ctx.Shutdown(stdcontext.TODO())
	err := ctx.ShutdownOrErr(stdcontext.TODO())
	suite.Equal("context lifecycle error: context already shutdown", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeLifecycle)
}

func (suite *LifecycleSuite) TestDependencyShutdownError() {
	foo1 := CtxAwareFoo{}
	foo2 := CtxAwareFoo{errOnShutdown: true}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctxb.Add(&foo2)
	ctx := ctxb.Build()
	err := ctx.ShutdownOrErr(stdcontext.TODO())
	suite.Equal("could not shutdown dependency: *di_test.CtxAwareFoo, cause:\nsimulated", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeDependencyShutdown)
}

func (suite *LifecycleSuite) TestInitAfterShutdownError() {
	foo1 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctx := ctxb.Build()
	ctx.Shutdown(stdcontext.TODO())
	err := ctx.InitializeOrErr()
	suite.Equal("context lifecycle error: context already shutdown", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeLifecycle)
}

func (suite *LifecycleSuite) TestGettingDependencyAfterShutdown() {
	foo1 := CtxAwareFoo{}
	ctxb := di.NewContextBuilder()
	ctxb.Add(&foo1)
	ctx := ctxb.Build()
	ctx.Shutdown(stdcontext.TODO())
	_, err := ctx.GetByTypeOrErr(new(Foo))
	suite.Equal("context lifecycle error: context already shutdown", err.Error())
	suite.Equal(err.ErrType(), di.ErrTypeLifecycle)
}

func TestLifecycleSuite(t *testing.T) {
	suite.Run(t, new(LifecycleSuite))
}
