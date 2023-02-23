package di_test

import (
	"fmt"
	"reflect"
	"testing"

	di "coditory.com/goiku-di"
	"github.com/stretchr/testify/suite"
)

type TypeValidationSuite struct {
	suite.Suite
}

func (suite *TypeValidationSuite) TestInvalidTypesOnAddAs() {
	tests := []struct {
		asType any
		value  any
	}{
		{
			asType: new(Bar),
			value:  &foo,
		},
		{
			asType: new(Bar),
			value:  foo,
		},
		{
			asType: new(int64),
			value:  42,
		},
	}

	for _, tt := range tests {
		desc := fmt.Sprintf("%s-%+v", reflect.TypeOf(tt.value), tt.value)
		suite.Run(desc, func() {
			ctxb := di.NewContextBuilder()
			err := func() (err error) {
				defer func() {
					err = recover().(error)
				}()
				ctxb.AddAs(tt.asType, tt.value)
				return nil
			}()
			suite.NotNil(err)
			suite.Equal(di.ErrInvalidType, err)
		})
	}
}

func TestTypeValidationSuite(t *testing.T) {
	suite.Run(t, new(TypeValidationSuite))
}
