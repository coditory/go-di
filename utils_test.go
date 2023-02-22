package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenericTypeOf(t *testing.T) {
	type Foo struct{}
	type Baz interface{}

	tests := []struct {
		gtype    func() reflect.Type
		expected string
	}{
		{
			gtype:    func() reflect.Type { return genericTypeOf[any]() },
			expected: "*interface {}",
		},
		{
			gtype:    func() reflect.Type { return genericTypeOf[*Foo]() },
			expected: "*di.Foo",
		},
		{
			gtype:    func() reflect.Type { return genericTypeOf[Foo]() },
			expected: "di.Foo",
		},
		{
			gtype:    func() reflect.Type { return genericTypeOf[Baz]() },
			expected: "*di.Baz",
		},
		{
			gtype:    func() reflect.Type { return genericTypeOf[[]Baz]() },
			expected: "[]di.Baz",
		},
	}

	for _, tt := range tests {
		desc := tt.expected
		t.Run(desc, func(t *testing.T) {
			actual := tt.gtype()
			assert.Equal(t, tt.expected, actual.String())
		})
	}
}
