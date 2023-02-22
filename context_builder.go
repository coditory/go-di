package di

import (
	"fmt"
	"reflect"
)

type ContextBuilder struct {
	holdersByCtors map[any]*holder
	holdersByType  map[reflect.Type][]*holder
	holdersByName  map[string]*holder
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		holdersByCtors: make(map[any]*holder),
		holdersByType:  make(map[reflect.Type][]*holder),
		holdersByName:  make(map[string]*holder),
	}
}

func (ctxb *ContextBuilder) Build() *Context {
	return &Context{
		holdersByType: ctxb.holdersByType,
		holdersByName: ctxb.holdersByName,
	}
}

func (ctxb *ContextBuilder) Add(ctor any) {
	holder, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	itype := holder.providesType
	ctxb.holdersByType[itype] = append(ctxb.holdersByType[itype], holder)
}

func (ctxb *ContextBuilder) AddAs(iface any, ctor any) {
	holder, err := createUniqueHolder(ctxb, ctor)
	if err != nil {
		panic(err)
	}
	itype := reflect.TypeOf(iface)
	ctxb.holdersByType[itype] = append(ctxb.holdersByType[itype], holder)
}

func createUniqueHolder(ctxb *ContextBuilder, ctor any) (*holder, error) {
	cval := reflect.ValueOf(ctor)
	if cval.Kind() != reflect.Func {
		return NewHolder(ctor)
	}
	ptr := fmt.Sprintf("%d", ctor)
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
