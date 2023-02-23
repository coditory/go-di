package di

import (
	"reflect"
)

func empty[T any]() (t T) {
	return
}

func genericTypeOf[T any]() reflect.Type {
	var t T
	ttype := reflect.TypeOf(t)
	if ttype == nil {
		ttype = reflect.TypeOf(new(T)).Elem()
	}
	return ttype
}
