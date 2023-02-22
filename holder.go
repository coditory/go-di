package di

import (
	"errors"
	"fmt"
	"reflect"
)

type Provider func(ctx *Context) (any, error)

type holder struct {
	provider     Provider
	created      bool
	instance     any
	providesType reflect.Type
}

func NewHolder(provider any) (*holder, error) {
	ptype := reflect.TypeOf(provider)
	if ptype == nil {
		return nil, errors.New("can't provide an untyped nil")
	}
	if ptype.Kind() == reflect.Func {
		return createLazyHolder(provider)
	} else {
		return createEagerHolder(provider)
	}
}

func createLazyHolder(provider any) (*holder, error) {
	pval := reflect.ValueOf(provider)
	ptype := pval.Type()
	numResults := ptype.NumOut()
	if numResults < 1 || numResults > 2 {
		return nil, fmt.Errorf("expected func to return one or two values")
	}
	if ptype.IsVariadic() {
		return nil, fmt.Errorf("variadic params not supported. use slice instead")
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if numResults == 2 && !ptype.Out(1).AssignableTo(errorInterface) {
		return nil, fmt.Errorf("expected func second result to be assignable to error")
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
		provider:     prov,
		providesType: resultType,
	}, nil
}

func createEagerHolder(value any) (*holder, error) {
	return &holder{
		created:      true,
		instance:     value,
		provider:     func(*Context) (any, error) { return value, nil },
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
				err = errors.New("provider panic")
			}
			result = empty[any]()
		}
	}()
	return holder.provider(ctx)
}
