package di

import (
	"errors"
	"reflect"
)

type ctor func(ctx *Context) (any, error)

type holder struct {
	ctor         ctor
	created      bool
	instance     any
	providesType reflect.Type
}

func newHolder(ctor any, lazy bool) (*holder, *Error) {
	ctype := reflect.TypeOf(ctor)
	if ctype == nil {
		return nil, newInvalidConstructorError("untyped constructor")
	}
	if lazy {
		return createLazyHolder(ctor)
	} else {
		return createEagerHolder(ctor)
	}
}

func createLazyHolder(ctor any) (*holder, *Error) {
	ctype := reflect.TypeOf(ctor)
	if ctype.Kind() != reflect.Func {
		return nil, newInvalidConstructorError("expected constructor function")
	}
	cval := reflect.ValueOf(ctor)
	numResults := ctype.NumOut()
	if numResults < 1 || numResults > 2 {
		return nil, newInvalidConstructorError("expected one result value with an optional error")
	}
	if ctype.IsVariadic() {
		return nil, newInvalidConstructorError("variadic parameters not supported (use slice instead)")
	}
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if numResults == 2 && !ctype.Out(1).AssignableTo(errorInterface) {
		return nil, newInvalidConstructorError("expected second result value to be an error")
	}
	resultType := ctype.Out(0)
	numArgs := ctype.NumIn()
	params := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		params[i] = ctype.In(i)
	}
	prov := func(ctx *Context) (any, error) {
		args := make([]reflect.Value, numArgs)
		for i, ptype := range params {
			if ptype.Kind() == reflect.Slice {
				argslice, err := ctx.getAllByRType(ptype.Elem())
				if err != nil {
					return nil, err
				}
				rslice := reflect.MakeSlice(ptype, 0, len(argslice))
				for _, item := range argslice {
					rslice = reflect.Append(rslice, reflect.ValueOf(item))
				}
				args[i] = rslice
			} else if ptype == genericTypeOf[*Context]() {
				args[i] = reflect.ValueOf(ctx)
			} else {
				arg, err := ctx.getByRType(ptype)
				if err != nil {
					return nil, err
				}
				args[i] = reflect.ValueOf(arg)
			}
		}
		result := cval.Call(args)
		obj := result[0].Interface()
		if len(result) == 2 {
			err, ok := result[1].Interface().(error)
			if !ok {
				return nil, newInvalidConstructorError("expected second result value to be an error")
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

func createEagerHolder(value any) (*holder, *Error) {
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
