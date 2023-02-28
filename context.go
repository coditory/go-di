package di

import (
	"errors"
	"fmt"
	"reflect"
)

type Context struct {
	path          map[string]int
	holdersByType map[reflect.Type][]*holder
	holdersByName map[string]*holder
}

func (ctx *Context) GetNamed(name string) (any, error) {
	holder := ctx.holdersByName[name]
	if holder == nil {
		return empty[any](), &MissingDependencyError{objName: &name}
	}
	depCtx, err := dependencyContext(ctx, descriptor(&name, nil))
	if err != nil {
		return empty[any](), err
	}
	obj, err := holder.getOrCreate(depCtx)
	if err != nil {
		if !errors.Is(err, ErrSkippedDependency) {
			creationErr := &DependencyCreationError{objName: &name, err: err}
			return empty[any](), creationErr
		}
	} else {
		return obj, nil
	}
	return empty[any](), &MissingDependencyError{objName: &name}
}

func (ctx *Context) getByType(otype reflect.Type) (any, error) {
	holders := ctx.holdersByType[otype]
	if holders == nil {
		return empty[any](), &MissingDependencyError{objType: &otype}
	}
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, descriptor(nil, &otype))
		if err != nil {
			return empty[any](), err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkippedDependency) {
				creationErr := &DependencyCreationError{objType: &otype, err: err}
				return empty[any](), creationErr
			}
		} else {
			return obj, nil
		}
	}
	return empty[any](), &MissingDependencyError{objType: &otype}
}

func (ctx *Context) GetAllByType(otype reflect.Type) ([]any, error) {
	holders := ctx.holdersByType[otype]
	result := make([]any, 0)
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, descriptor(nil, &otype))
		if err != nil {
			return nil, err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkippedDependency) {
				creationErr := &DependencyCreationError{objType: &otype, err: err}
				return nil, creationErr
			}
		} else {
			result = append(result, obj)
		}
	}
	return result, nil
}

func GetOrPanic[T any](ctx *Context) T {
	obj, err := Get[T](ctx)
	if err != nil {
		panic(err)
	}
	return obj
}

func Get[T any](ctx *Context) (T, error) {
	ttype := genericTypeOf[T]()
	obj, err := ctx.getByType(ttype)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), &InvalidTypeError{objType: reflect.TypeOf(obj), expectedType: genericTypeOf[T]()}
	}
	return typed, nil
}

func GetNamedOrPanic[T any](ctx *Context, name string) T {
	obj, err := GetNamed[T](ctx, name)
	if err != nil {
		panic(err)
	}
	return obj
}

func GetNamed[T any](ctx *Context, name string) (T, error) {
	obj, err := ctx.GetNamed(name)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), &InvalidTypeError{objName: &name, objType: reflect.TypeOf(obj), expectedType: genericTypeOf[T]()}
	}
	return typed, nil
}

func GetAllOrPanic[T any](ctx *Context) []T {
	result, err := GetAll[T](ctx)
	if err != nil {
		panic(err)
	}
	return result
}

func GetAll[T any](ctx *Context) ([]T, error) {
	ttype := genericTypeOf[T]()
	objs, err := ctx.GetAllByType(ttype)
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

func dependencyContext(ctx *Context, descriptor string) (*Context, error) {
	if ctx.path[descriptor] > 0 {
		cycle := make([]string, len(ctx.path))
		for d, i := range ctx.path {
			cycle[i-1] = d
		}
		return nil, &CyclicDependencyError{path: cycle}
	}
	path := make(map[string]int)
	for k, v := range ctx.path {
		path[k] = v
	}
	path[descriptor] = len(path) + 1
	sub := Context{
		path:          path,
		holdersByType: ctx.holdersByType,
		holdersByName: ctx.holdersByName,
	}
	return &sub, nil
}

func descriptor(objName *string, objType *reflect.Type) string {
	var result string
	if objName != nil && objType != nil {
		result = fmt.Sprintf("%s (name: %s)", (*objType).String(), *objName)
	} else if objName != nil {
		return fmt.Sprintf("(name: %s)", *objName)
	} else if objType != nil {
		return (*objType).String()
	} else {
		panic("Expected obj name or obj type to be defined")
	}
	return result
}
