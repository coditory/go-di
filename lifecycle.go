package di

import (
	stdcontext "context"
	"reflect"
)

type Shutdownable interface {
	Shutdown(context stdcontext.Context)
}

type Initializable interface {
	Initialize()
}

var (
	initializableType  = new(Initializable)
	initializableRType = reflect.TypeOf(initializableType).Elem()
	shutdownableType   = new(Shutdownable)
	shutdownableRType  = reflect.TypeOf(shutdownableType).Elem()
)
