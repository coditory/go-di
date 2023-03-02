package di_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type DependencyValidationSuite struct {
	suite.Suite
}

func (suite *DependencyValidationSuite) TestInvalidTypesOnAddAs() {
	tests := []struct {
		asType  any
		value   any
		message string
	}{
		{
			asType:  new(Bar),
			value:   &foo,
			message: "could not cast *di_test.Foo to di_test.Bar",
		},
		{
			asType:  new(Bar),
			value:   foo,
			message: "could not cast di_test.Foo to di_test.Bar",
		},
		{
			asType:  new(Foo),
			value:   &foo,
			message: "could not cast *di_test.Foo to di_test.Foo",
		},
		{
			asType:  new(*Foo),
			value:   foo,
			message: "could not cast di_test.Foo to *di_test.Foo",
		},
		{
			asType:  new(int64),
			value:   42,
			message: "could not cast int to int64",
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s %+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			err := ctxb.AddAsOrErr(tt.asType, tt.value)
			suite.NotNil(err)
			suite.Equal(tt.message, err.Error())
			suite.Equal(di.ErrTypeInvalidType, err.ErrType())
		})
	}
}

func TestDependencyValidationSuite(t *testing.T) {
	suite.Run(t, new(DependencyValidationSuite))
}
