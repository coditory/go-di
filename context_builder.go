package di

import (
	"fmt"
	"reflect"
)

type ContextBuilder struct {
	holdersByCtors map[any]*holder
	holdersByType  map[reflect.Type]*Set[*holder]
	holdersByName  map[string]*holder
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		holdersByCtors: make(map[any]*holder),
		holdersByType:  make(map[reflect.Type]*Set[*holder]),
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
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	err = ctxb.addHolderForType(hldr, hldr.providesType)
	if err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddNamed(name string, ctor any) {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	err = ctxb.addHolderForName(hldr, name)
	if err != nil {
		panic(err)
	}
	err = ctxb.addHolderForType(hldr, hldr.providesType)
	if err != nil {
		delete(ctxb.holdersByName, name)
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddAs(iface any, ctor any) {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	itype := reflect.TypeOf(iface).Elem()
	err = ctxb.addHolderForType(hldr, itype)
	if err != nil {
		panic(err)
	}
}

func (ctxb *ContextBuilder) AddNamedAs(name string, iface any, ctor any) {
	hldr, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	err = ctxb.addHolderForName(hldr, name)
	if err != nil {
		panic(err)
	}
	itype := reflect.TypeOf(iface).Elem()
	err = ctxb.addHolderForType(hldr, itype)
	if err != nil {
		delete(ctxb.holdersByName, name)
		panic(err)
	}
}

func (ctxb *ContextBuilder) addHolderForType(hldr *holder, itype reflect.Type) error {
	if hldr.providesType != itype && !hldr.providesType.AssignableTo(itype) {
		return &InvalidTypeError{objType: hldr.providesType, expectedType: itype}
	}
	if ctxb.holdersByType[itype] == nil {
		ctxb.holdersByType[itype] = NewSet[*holder]()
	}
	if ctxb.holdersByType[itype].Contains(hldr) {
		return ErrDuplicatedRegistration
	}
	ctxb.holdersByType[itype].Add(hldr)
	return nil
}

func (ctxb *ContextBuilder) addHolderForName(hldr *holder, name string) error {
	if ctxb.holdersByName[name] != nil {
		return &DuplicatedNameError{name: name}
	}
	ctxb.holdersByName[name] = hldr
	return nil
}

func createUniqueHolder(ctxb *ContextBuilder, ctor any) (*holder, error) {
	cval := reflect.ValueOf(ctor)
	ckind := cval.Kind()
	var ptr string
	if ckind == reflect.Pointer && cval.IsNil() {
		ptr = fmt.Sprintf("nil-%T", ctor)
	} else if ckind == reflect.Func || ckind == reflect.Pointer {
		ptr = fmt.Sprintf("ptr-%p", ctor)
	} else {
		return NewHolder(ctor)
	}
	hldr := ctxb.holdersByCtors[ptr]
	if hldr == nil {
		nhldr, err := NewHolder(ctor)
		if err != nil {
			return nil, err
		}
		ctxb.holdersByCtors[ptr] = nhldr
		hldr = nhldr
	}
	return hldr, nil
}
