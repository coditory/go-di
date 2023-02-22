package di_test

type Baz interface{ baz() }

type Foo struct{}

func (f *Foo) baz() {}

type Bar struct{}

func (b Bar) baz() {}

var (
	foo Foo = Foo{}
	bar Bar = Bar{}
)
