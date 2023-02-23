package di_test

import (
	"fmt"
	"reflect"
	"testing"

	di "github.com/coditory/go-di"
	"github.com/stretchr/testify/suite"
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
			err := func() (err any) {
				defer func() {
					err = recover()
				}()
				ctxb.Add(tt.value)
				ctxb.Add(tt.value)
				return nil
			}()
			suite.NotNil(err)
			suite.Equal("duplicated registration", err.(error).Error())
		})
	}
}

func (suite *DuplicatedRegistartionSuite) TestAllowDuplicatedNonPointers() {
	tests := []struct {
		value any
	}{
		{value: make([]string, 1)},
		{value: []string{"abc"}},
		{value: make(map[int]string)},
		{value: "abc"},
		{value: 42},
		{value: foo},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			err := func() (err any) {
				defer func() {
					err = recover()
				}()
				ctxb.Add(tt.value)
				ctxb.Add(tt.value)
				return nil
			}()
			suite.Nil(err)
			ctx := ctxb.Build()
			objs, err := ctx.GetAllByType(reflect.TypeOf(tt.value))
			suite.Nil(err)
			suite.Equal(2, len(objs))
			suite.Equal(tt.value, objs[0])
			suite.Equal(tt.value, objs[1])
		})
	}
}

func TestDuplicatedRegistartionSuite(t *testing.T) {
	suite.Run(t, new(DuplicatedRegistartionSuite))
}
