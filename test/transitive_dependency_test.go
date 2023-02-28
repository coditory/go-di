package di_test

import (
	"testing"

	di "github.com/coditory/go-di"
	"github.com/stretchr/testify/assert"
)

type transitiveFoo struct {
	bar *transitiveBar
}

type transitiveBar struct {
	baz *transitiveBaz
}

type transitiveBaz struct {
	message string
}

func provideTransitiveFoo(bar *transitiveBar) *transitiveFoo {
	return &transitiveFoo{bar: bar}
}

func provideTransitiveBar(baz *transitiveBaz) *transitiveBar {
	return &transitiveBar{baz: baz}
}

func provideTransitiveBaz() *transitiveBaz {
	return &transitiveBaz{message: "hello"}
}

func Test_TransitiveDependency(t *testing.T) {
	ctxb := di.NewContextBuilder()
	ctxb.Add(provideTransitiveFoo)
	ctxb.Add(provideTransitiveBar)
	ctxb.Add(provideTransitiveBaz)
	ctx := ctxb.Build()
	result, err := di.GetOrErr[*transitiveFoo](ctx)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.bar)
	assert.NotNil(t, result.bar.baz)
	assert.Equal(t, result.bar.baz.message, "hello")
}
