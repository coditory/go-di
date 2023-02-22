package di

import (
	"reflect"
)

type ContextBuilder struct {
	holdersByProviders map[any]*holder
	holdersByType      map[reflect.Type][]*holder
	holdersByName      map[string]*holder
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		holdersByProviders: make(map[any]*holder),
		holdersByType:      make(map[reflect.Type][]*holder),
		holdersByName:      make(map[string]*holder),
	}
}

func (ctxb *ContextBuilder) Build() *Context {
	return &Context{
		holdersByType: ctxb.holdersByType,
		holdersByName: ctxb.holdersByName,
	}
}

func (ctxb *ContextBuilder) Add(provider any) {
	holder, err := createUniqueHolder(ctxb, provider)
	if err != nil {
		panic(err)
	}
	itype := holder.providesType
	ctxb.holdersByType[itype] = append(ctxb.holdersByType[itype], holder)
}

func (ctxb *ContextBuilder) AddAs(iface any, provider any) {
	holder, err := createUniqueHolder(ctxb, provider)
	if err != nil {
		panic(err)
	}
	itype := reflect.TypeOf(iface)
	ctxb.holdersByType[itype] = append(ctxb.holdersByType[itype], holder)
}

func createUniqueHolder(ctxb *ContextBuilder, provider any) (*holder, error) {
	ptr, ok := provider.(*any)
	if !ok {
		ptr = &provider
	}
	hldr := ctxb.holdersByProviders[ptr]
	if hldr == nil {
		nhldr, err := NewHolder(provider)
		if err != nil {
			return nil, err
		}
		ctxb.holdersByProviders[ptr] = nhldr
		hldr = nhldr
	}
	return hldr, nil
}
