package di

import (
	"reflect"
)

func GetOrPanic[T any](ctx *Context) T {
	obj, err := GetOrErr[T](ctx)
	if err != nil {
		panic(err)
	}
	return obj
}

func Get[T any](ctx *Context) T {
	obj, err := GetOrErr[T](ctx)
	if err != nil {
		panic(err)
	}
	return obj
}

func GetOrErr[T any](ctx *Context) (T, error) {
	ttype := genericTypeOf[T]()
	obj, err := ctx.getByRType(ttype)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), &InvalidTypeError{objType: reflect.TypeOf(obj), expectedType: genericTypeOf[T]()}
	}
	return typed, nil
}

func GetNamed[T any](ctx *Context, name string) T {
	obj, err := GetNamedOrErr[T](ctx, name)
	if err != nil {
		panic(err)
	}
	return obj
}

func GetNamedOrErr[T any](ctx *Context, name string) (T, error) {
	obj, err := ctx.GetNamedOrErr(name)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), &InvalidTypeError{objName: &name, objType: reflect.TypeOf(obj), expectedType: genericTypeOf[T]()}
	}
	return typed, nil
}

func GetAll[T any](ctx *Context) []T {
	result, err := GetAllOrErr[T](ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func GetAllOrErr[T any](ctx *Context) ([]T, error) {
	ttype := genericTypeOf[T]()
	objs, err := ctx.getAllByRType(ttype)
	if err != nil {
		return nil, err
	}
	result := make([]T, 0)
	for _, obj := range objs {
		if typed, ok := obj.(T); ok {
			result = append(result, typed)
		} else {
			return nil, &InvalidTypeError{objType: ttype, expectedType: genericTypeOf[T]()}
		}
	}
	return result, nil
}
