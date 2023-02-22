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
		ttype = reflect.TypeOf(new(T))
	}
	// fmt.Println("ttype:", ttype)
	// fmt.Println("ttype nil:", ttype == nil)
	// fmt.Println("ttype kind:", ttype.Kind())
	// fmt.Println("ttype value:", reflect.ValueOf(t))
	// fmt.Println("ttype Y:", ttype)
	return ttype
}
