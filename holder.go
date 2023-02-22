package di

import (
	"errors"
	"fmt"
	"reflect"
)

type ctor func(ctx *Context) (any, error)

type holder struct {
	ctor         ctor
	created      bool
	instance     any
	providesType reflect.Type
}

func NewHolder(ctor any) (*holder, error) {
	ctype := reflect.TypeOf(ctor)
	if ctype == nil {
		return nil, errors.New("untyped constructor")
	}
	if ctype.Kind() == reflect.Func {
		return createLazyHolder(ctor)
	} else {
		return createEagerHolder(ctor)
	}
}

func createLazyHolder(ctor any) (*holder, error) {
	pval := reflect.ValueOf(ctor)
	ptype := pval.Type()
	numResults := ptype.NumOut()
	if numResults < 1 || numResults > 2 {
		return nil, fmt.Errorf("expected constructor to return one result")
	}
	if ptype.IsVariadic() {
		return nil, fmt.Errorf("variadic constructor parameters are not supported (use slice instead)")
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if numResults == 2 && !ptype.Out(1).AssignableTo(errorInterface) {
		return nil, fmt.Errorf("expected constructor second result to be an error")
	}
	resultType := ptype.Out(0)
	numArgs := ptype.NumIn()
	params := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		params[i] = ptype.In(i)
	}
	prov := func(ctx *Context) (any, error) {
		args := make([]reflect.Value, numArgs)
		for i, paramType := range params {
			if paramType.Kind() == reflect.Slice {
				argslice, err := ctx.getAllByType(paramType)
				if err != nil {
					return nil, err
				}
				args[i] = reflect.ValueOf(argslice)
			} else {
				arg, err := ctx.getByType(paramType)
				if err != nil {
					return nil, err
				}
				args[i] = reflect.ValueOf(arg)
			}
		}
		result := pval.Call(args)
		obj := result[0].Interface()
		if len(result) == 2 {
			err, ok := result[1].Interface().(error)
			if !ok {
				return nil, ErrInvalidType
			}
			return obj, err
		}
		return obj, nil
	}
	return &holder{
		ctor:         prov,
		providesType: resultType,
	}, nil
}

func createEagerHolder(value any) (*holder, error) {
	return &holder{
		created:      true,
		instance:     value,
		ctor:         func(*Context) (any, error) { return value, nil },
		providesType: reflect.TypeOf(value),
	}, nil
}

func (h *holder) getOrCreate(ctx *Context) (any, error) {
	obj := h.instance
	if !h.created {
		newobj, err := provide(ctx, h)
		if err != nil {
			return empty[any](), err
		}
		h.instance = newobj
		h.created = true
		obj = newobj
	}
	return obj, nil
}

func provide(ctx *Context, holder *holder) (result any, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("ctor panic")
			}
			result = empty[any]()
		}
	}()
	return holder.ctor(ctx)
}
