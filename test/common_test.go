package di_test

import "errors"

type Baz interface {
	Id() string
}

type Foo struct {
	id string
}

func (f *Foo) Id() string {
	return f.id
}

type Bar struct {
	id string
}

func (b Bar) Id() string {
	return b.id
}

var (
	foo          Foo   = Foo{id: "foo"}
	foo2         Foo   = Foo{id: "foo2"}
	bar          Bar   = Bar{id: "bar"}
	errSimulated error = errors.New("simulated")
)
