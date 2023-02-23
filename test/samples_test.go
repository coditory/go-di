package di_test

import "errors"

type Baz interface{ baz() }

type Foo struct {
	name string
}

func (f *Foo) baz() {}

type Bar struct {
	name string
}

func (b Bar) baz() {}

var (
	foo          Foo   = Foo{name: "foo"}
	bar          Bar   = Bar{name: "bar"}
	errSimulated error = errors.New("simulated")
)
