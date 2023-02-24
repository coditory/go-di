package di

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrMissingObject          = errors.New("missing object")
	ErrInvalidType            = errors.New("invalid type")
	ErrDuplicatedName         = errors.New("duplicated name")
	ErrDuplicatedRegistration = errors.New("duplicated registration")
	ErrDependencyCycle        = errors.New("dependency cycle")
	ErrSkipped                = errors.New("registration skipped")
)

type Context struct {
	path          map[*holder]int
	holdersByType map[reflect.Type][]*holder
	holdersByName map[string]*holder
}

func (ctx *Context) GetNamed(name string) (any, error) {
	holder := ctx.holdersByName[name]
	if holder == nil {
		return empty[any](), ErrMissingObject
	}
	depCtx, err := dependencyContext(ctx, holder)
	if err != nil {
		return empty[any](), err
	}
	obj, err := holder.getOrCreate(depCtx)
	if err != nil {
		if !errors.Is(err, ErrSkipped) {
			return empty[any](), err
		}
	} else {
		return obj, nil
	}
	return empty[any](), ErrMissingObject
}

func (ctx *Context) GetByType(otype reflect.Type) (any, error) {
	holders := ctx.holdersByType[otype]
	if holders == nil {
		return empty[any](), ErrMissingObject
	}
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, holder)
		if err != nil {
			return empty[any](), err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkipped) {
				return empty[any](), err
			}
		} else {
			return obj, nil
		}
	}
	return empty[any](), ErrMissingObject
}

func (ctx *Context) GetAllByType(otype reflect.Type) ([]any, error) {
	holders := ctx.holdersByType[otype]
	result := make([]any, 0)
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, holder)
		if err != nil {
			return nil, err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkipped) {
				return nil, err
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
	obj, err := ctx.GetByType(ttype)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), fmt.Errorf("invalid type - could not cast %T to %s", obj, ttype)
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
	ttype := genericTypeOf[T]()
	obj, err := ctx.GetNamed(name)
	if err != nil {
		return empty[T](), err
	}
	typed, ok := obj.(T)
	if !ok {
		return empty[T](), fmt.Errorf("invalid type - could not cast %T to %s", obj, ttype)
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
			return nil, ErrInvalidType
		}
	}
	return result, nil
}

func dependencyContext(ctx *Context, hldr *holder) (*Context, error) {
	if ctx.path[hldr] > 0 {
		return nil, ErrDependencyCycle
	}
	path := make(map[*holder]int)
	for k, v := range ctx.path {
		path[k] = v
	}
	path[hldr] = len(path) + 1
	sub := Context{
		path:          path,
		holdersByType: ctx.holdersByType,
		holdersByName: ctx.holdersByName,
	}
	return &sub, nil
}
