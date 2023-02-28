package di_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type DuplicatedRegistartionSuite struct {
	suite.Suite
}

func (suite *DuplicatedRegistartionSuite) TestWithFunction() {
	inits := 0
	ctor := func() *Foo { inits++; return &Foo{} }
	ctxb := di.NewContextBuilder()
	err := func() (err any) {
		defer func() {
			err = recover()
		}()
		ctxb.Add(ctor)
		ctxb.Add(ctor)
		return nil
	}()
	suite.NotNil(err)
	suite.Equal("duplicated registration", err.(error).Error())
	suite.Equal(0, inits)
}

func (suite *DuplicatedRegistartionSuite) TestForbiddenDuplicatedPointers() {
	slice := make([]string, 1)
	text := "abc"
	tests := []struct {
		value any
	}{
		{value: &slice},
		{value: &text},
		{value: &foo},
		{value: (*Foo)(nil)},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.Add(tt.value)
			err := ctxb.AddOrErr(tt.value)
			suite.NotNil(err)
			suite.Equal("duplicated registration", err.Error())
		})
	}
}

func (suite *DuplicatedRegistartionSuite) TestAllowDuplicatedNonPointers() {
	tests := []struct {
		value any
		atype any
	}{
		{value: make([]string, 1), atype: new([]string)},
		{value: []string{"abc"}, atype: new([]string)},
		{value: make(map[int]string), atype: new(map[int]string)},
		{value: "abc", atype: new(string)},
		{value: 42, atype: new(int)},
		{value: foo, atype: new(Foo)},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			ctxb.Add(tt.value)
			ctxb.Add(tt.value)
			ctx := ctxb.Build()
			objs := ctx.GetAllByType(tt.atype)
			suite.Equal(2, len(objs))
			suite.Equal(tt.value, objs[0])
			suite.Equal(tt.value, objs[1])
		})
	}
}

func TestDuplicatedRegistartionSuite(t *testing.T) {
	suite.Run(t, new(DuplicatedRegistartionSuite))
}
