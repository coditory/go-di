package di

import (
	"fmt"
	"reflect"

	coll "github.com/coditory/go-di/internal/collection"
)

type ContextBuilder struct {
	holdersByCtors map[any]*holder
	holdersByType  map[reflect.Type]*coll.Set[*holder]
	holdersByName  map[string]*holder
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		holdersByCtors: make(map[any]*holder),
		holdersByType:  make(map[reflect.Type]*coll.Set[*holder]),
		holdersByName:  make(map[string]*holder),
	}
}

func (ctxb *ContextBuilder) Build() *Context {
	holders := make(map[reflect.Type][]*holder)
	for k, v := range ctxb.holdersByType {
		holders[k] = v.ToSlice()
	}
	return &Context{
		holdersByType: holders,
		holdersByName: ctxb.holdersByName,
	}
}

func (ctxb *ContextBuilder) Add(ctor any) {
	if err := ctxb.AddOrErr(ctor); err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddOrErr(ctor any) *Error {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		return err
	}
	err = ctxb.addHolderForType(hldr, hldr.providesType)
	if err != nil {
		return err
	}
	return nil
}

func (ctxb *ContextBuilder) AddNamed(name string, ctor any) {
	if err := ctxb.AddNamedOrErr(name, ctor); err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddNamedOrErr(name string, ctor any) *Error {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		return err
	}
	err = ctxb.addHolderForName(hldr, name)
	if err != nil {
		return err
	}
	err = ctxb.addHolderForType(hldr, hldr.providesType)
	if err != nil {
		delete(ctxb.holdersByName, name)
		return err
	}
	return nil
}

func (ctxb *ContextBuilder) AddAs(atype any, ctor any) {
	if err := ctxb.AddAsOrErr(atype, ctor); err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddAsOrErr(atype any, ctor any) *Error {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		return err
	}
	rtype := reflect.TypeOf(atype).Elem()
	err = ctxb.addHolderForType(hldr, rtype)
	if err != nil {
		return err
	}
	return nil
}

func (ctxb *ContextBuilder) AddNamedAs(name string, atype any, ctor any) {
	if err := ctxb.AddNamedAsOrErr(name, atype, ctor); err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddNamedAsOrErr(name string, atype any, ctor any) *Error {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		return err
	}
	err = ctxb.addHolderForName(hldr, name)
	if err != nil {
		return err
	}
	rtype := reflect.TypeOf(atype).Elem()
	err = ctxb.addHolderForType(hldr, rtype)
	if err != nil {
		delete(ctxb.holdersByName, name)
		return err
	}
	return nil
}

func (ctxb *ContextBuilder) addHolderForType(hldr *holder, rtype reflect.Type) *Error {
	if hldr.providesType != rtype && !hldr.providesType.AssignableTo(rtype) {
		return newInvalidTypeError(nil, hldr.providesType, rtype)
	}
	if ctxb.holdersByType[rtype] == nil {
		ctxb.holdersByType[rtype] = coll.NewSet[*holder]()
	}
	if ctxb.holdersByType[rtype].Contains(hldr) {
		return newDuplicatedRegistrationError()
	}
	ctxb.holdersByType[rtype].Add(hldr)
	return nil
}

func (ctxb *ContextBuilder) addHolderForName(hldr *holder, name string) *Error {
	if ctxb.holdersByName[name] != nil {
		return newDuplicatedNameError(name)
	}
	ctxb.holdersByName[name] = hldr
	return nil
}

func createUniqueHolder(ctxb *ContextBuilder, ctor any) (*holder, *Error) {
	cval := reflect.ValueOf(ctor)
	ckind := cval.Kind()
	var ptr string
	if ckind == reflect.Pointer && cval.IsNil() {
		ptr = fmt.Sprintf("nil-%T", ctor)
	} else if ckind == reflect.Func || ckind == reflect.Pointer {
		ptr = fmt.Sprintf("ptr-%p", ctor)
	} else {
		return newHolder(ctor)
	}
	hldr := ctxb.holdersByCtors[ptr]
	if hldr == nil {
		nhldr, err := newHolder(ctor)
		if err != nil {
			return nil, err
		}
		ctxb.holdersByCtors[ptr] = nhldr
		hldr = nhldr
	}
	return hldr, nil
}
