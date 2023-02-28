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

func (ctx *Context) GetNamed(name string) any {
	obj, err := ctx.GetNamedOrErr(name)
	if err != nil {
		panic(err)
	}
	return obj
}

func (ctx *Context) GetNamedOrErr(name string) (any, error) {
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

func (ctx *Context) GetByType(atype any) any {
	obj, err := ctx.GetByTypeOrErr(atype)
	if err != nil {
		panic(err)
	}
	return obj
}

func (ctx *Context) GetByTypeOrErr(atype any) (any, error) {
	rtype := reflect.TypeOf(atype).Elem()
	return ctx.getByRType(rtype)
}

func (ctx *Context) GetAllByType(atype any) []any {
	obj, err := ctx.GetAllByTypeOrErr(atype)
	if err != nil {
		panic(err)
	}
	return obj
}

func (ctx *Context) GetAllByTypeOrErr(atype any) ([]any, error) {
	rtype := reflect.TypeOf(atype).Elem()
	return ctx.getAllByRType(rtype)
}

func (ctx *Context) getByRType(rtype reflect.Type) (any, error) {
	holders := ctx.holdersByType[rtype]
	if holders == nil {
		return empty[any](), &MissingDependencyError{objType: &rtype}
	}
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, descriptor(nil, &rtype))
		if err != nil {
			return empty[any](), err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkippedDependency) {
				creationErr := &DependencyCreationError{objType: &rtype, err: err}
				return empty[any](), creationErr
			}
		} else {
			return obj, nil
		}
	}
	return empty[any](), &MissingDependencyError{objType: &rtype}
}

func (ctx *Context) getAllByRType(rtype reflect.Type) ([]any, error) {
	holders := ctx.holdersByType[rtype]
	result := make([]any, 0)
	for _, holder := range holders {
		depCtx, err := dependencyContext(ctx, descriptor(nil, &rtype))
		if err != nil {
			return nil, err
		}
		obj, err := holder.getOrCreate(depCtx)
		if err != nil {
			if !errors.Is(err, ErrSkippedDependency) {
				creationErr := &DependencyCreationError{objType: &rtype, err: err}
				return nil, creationErr
			}
		} else {
			result = append(result, obj)
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
