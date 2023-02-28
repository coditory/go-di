package di_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	di "github.com/coditory/go-di"
)

type CreatorValidationSuite struct {
	suite.Suite
}

func (suite *CreatorValidationSuite) TestErrorResult() {
	tests := []struct {
		title string
		ctor  any
		error string
	}{
		{
			title: "two results",
			error: "invalid dependency constructor: expected second result value to be an error",
			ctor: func() (*Foo, *Bar) {
				return nil, nil
			},
		},
		{
			title: "three results, last err",
			error: "invalid dependency constructor: expected one result value with an optional error",
			ctor: func() (*Foo, *Bar, error) {
				return nil, nil, nil
			},
		},
		{
			title: "zero results",
			error: "invalid dependency constructor: expected one result value with an optional error",
			ctor: func() {
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.title, func() {
			ctxb := di.NewContextBuilder()
			err := func() (err error) {
				defer func() {
					err = recover().(error)
				}()
				ctxb.Add(tt.ctor)
				return nil
			}()
			suite.Equal(tt.error, err.Error())
		})
	}
}

func TestCreatorValidationSuite(t *testing.T) {
	suite.Run(t, new(CreatorValidationSuite))
}
