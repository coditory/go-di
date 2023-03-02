# Coditory - Go Dependency Injection
[![GitHub release](https://img.shields.io/github/v/release/coditory/go-di.svg)](https://github.com/coditory/go-di/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/coditory/go-di.svg)](https://pkg.go.dev/github.com/coditory/go-di)
[![Go Report Card](https://goreportcard.com/badge/github.com/coditory/go-di)](https://goreportcard.com/report/github.com/coditory/go-di)
[![Build Status](https://github.com/coditory/go-di/workflows/Build/badge.svg?branch=main)](https://github.com/coditory/go-di/actions?query=workflow%3ABuild+branch%3Amain)
[![Coverage](https://codecov.io/gh/coditory/go-di/branch/main/graph/badge.svg?token=EPRs5LiPje)](https://codecov.io/gh/coditory/go-di)

**🚧 This library as under heavy development until release of version `1.x.x` 🚧**

> Dependency injection for Go projects, targeted for web applications that create DI context during the startup and operates on collections of dependencies.

- Register dependencies by name and type
- Register multiple dependencies per type
- Conditional dependency registration
- Simple dependency retrieval - no manual casting or additional callbacks
- Simple setup - no generators
- Detection of slow dependency creation (TODO)
- Initialization and finalization mechanisms (TODO)

# Getting started

## Installation
Get the dependency with:
```sh
go get github.com/coditory/go-di
```

and import it in the project:
```go
import "github.com/coditory/go-di"
```

The exported package is `di`, basic usage:
```go
import "github.com/coditory/go-di"

func main() {
  ctxb := di.NewContextBuilder()
  // Add dependencies
  ctxb.Add(&foo)
  // Build the di context
  ctx := ctxb.Build()
  // Retrieve dependencies from the context
}
```

## Full example

```go
package main

import (
  "fmt"
  "github.com/coditory/go-di"
)

type Baz interface { Id() string }
type Foo struct { id string }
func (f *Foo) Id() string { return f.id }
type Bar struct { id string }
func (b *Bar) Id() string { return b.id }

func main() {
  foo := Foo{id: "foo1"}
  foo2 := Foo{id: "foo2"}
  bar := Bar{id: "baz"}

  ctxb := di.NewContextBuilder()
  ctxb.AddAs(new(Baz), &foo)
  ctxb.AddAs(new(Baz), &foo2)
  ctxb.AddAs(new(Baz), &bar)
  ctx := ctxb.Build()

  fmt.Printf("first implementation: %+v\n", di.Get[Baz](ctx))
  for i, baz := range di.GetAll[Baz](ctx) {
    fmt.Printf("%d: %+v\n", i, baz)
  }
}
// output:
// first implementation: &{id:foo1}
// 0: &{id:foo1}
// 1: &{id:foo2}
// 2: &{id:baz}
```

# Usage

## Dependencies by type

Add dependency reference and retrieve by the reference type:
```go
ctxb := di.NewContextBuilder()
ctxb.Add(&foo)
ctxb.Add(&foo2)
ctx := ctxb.Build()
suite.Equal(&foo, di.Get[*Foo](ctx))
suite.Equal([]*Foo{&foo, &foo2}, di.GetAll[*Foo](ctx))
```

Add dependency (by value) and retrieve by the type:
```go
ctxb := di.NewContextBuilder()
ctxb.Add(foo)
ctxb.Add(foo2)
ctx := ctxb.Build()
suite.Equal(foo, di.Get[Foo](ctx))
foos := di.GetAll[Foo](ctx)
suite.Equal(2, len(foos))
suite.Equal("foo", foos[0].Id())
suite.Equal("foo2", foos[1].Id())
```

## Dependencies by interface

To retrieve a dependency by the interface it must be registered explicitly by the interface type:
```go
ctxb := di.NewContextBuilder()
// ctx.Add(&foo) <- will not work!
ctxb.AddAs(new(Baz), &foo)
ctxb.AddAs(new(Baz), &foo2)
ctx := ctxb.Build()
suite.Equal(&foo, di.Get[Baz](ctx))
suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
```

Single dependency can be registered multiple times:
```go
ctxb := di.NewContextBuilder()
ctxb.Add(&foo)
ctxb.AddAs(new(Baz), &foo)
ctxb.Add(&foo2)
ctxb.AddAs(new(Baz), &foo2)
ctx := ctxb.Build()
suite.Equal(&foo, di.Get[Baz](ctx))
suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
suite.Equal(&foo, di.Get[*Foo](ctx))
suite.Equal([]*Foo{&foo, &foo2}, di.GetAll[*Foo](ctx))
```

## Named dependencies

To differentiate dependencies of the same type you can name them.
There can be only one dependency per name.

```go
ctxb := di.NewContextBuilder()
ctxb.AddAs(new(Baz), &foo)
ctxb.AddNamedAs("special-foo", new(Baz), &foo2)
// ctxb.AddNamedAs("special-foo", new(Baz), &foo3) <- panics
ctx := ctxb.Build()
suite.Equal(&foo, di.Get[Baz](ctx))
suite.Equal(&foo2, di.GetNamed[Baz](ctx, "special-foo"))
suite.Equal([]Baz{&foo, &foo2}, di.GetAll[Baz](ctx))
```

## Lazy dependencies

Lazy dependencies are created when retrieved:

```go
type Boo struct {
  foo *Foo
  bar *Bar
}
ctxb := di.NewContextBuilder()
ctxb.Provide(func() *Foo {
  return &foo
})
ctxb.Provide(func() *Bar {
  return &bar
})
ctxb.Provide(func(foo *Foo, bar *Bar) *Boo {
  return &Boo{foo: foo, bar: bar}
})
ctx := ctxb.Build()
boo := di.Get[*Boo](ctx)
suite.Equal(&foo, boo.foo)
suite.Equal(&bar, boo.bar)
```

If a dependency requires a long list of other dependencies then inject `*Context`:
```go
ctxb.Provide(func(ctx *di.Context) *Boo {
  return &Boo{foo: di.Get[*Foo](ctx), bar: di.Get[*Bar](ctx)}
})
```

There are the same variations for `ctxb.Add*` and `ctxb.Provide*` functions:
```go
ctxb.Provide(createFoo)
ctxb.ProvideAs(new(Baz), createFoo)
ctxb.ProvideNamed("special-foo", createFoo)
ctxb.ProvideNamedAs("special-foo", new(Baz), createFoo)
```

When dependency creator is a separate (non-inlined) function, then creation happens only once:
```go
creations := 0
createFoo := func() *Foo {
  creations++
  return &foo
}
ctxb := di.NewContextBuilder()
ctxb.Provide(createFoo)
ctxb.ProvideAs(new(Baz), createFoo)
ctx := ctxb.Build()
di.Get[Baz](ctx)
di.Get[*Foo](ctx)
suite.Equal(1, creations)
```
