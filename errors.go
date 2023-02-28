package di

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrSkippedDependency      = errors.New("skipped dependency")
	ErrDuplicatedRegistration = errors.New("duplicated registration")
)

type DependencyCreationError struct {
	objName *string
	objType *reflect.Type
	err     error
}

func (e *DependencyCreationError) Error() string {
	return fmt.Sprintf("could not create dependency %s, cause:\n%s", descriptor(e.objName, e.objType), e.err)
}

func (e *DependencyCreationError) Unwrap() error {
	return e.err
}

type MissingDependencyError struct {
	objName *string
	objType *reflect.Type
}

func (e *MissingDependencyError) Error() string {
	return fmt.Sprintf("missing dependency %s", descriptor(e.objName, e.objType))
}

type InvalidTypeError struct {
	objName      *string
	objType      reflect.Type
	expectedType reflect.Type
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("could not cast %s to %s", descriptor(e.objName, &e.objType), e.expectedType)
}

type CyclicDependencyError struct {
	path []string
}

func (e *CyclicDependencyError) Error() string {
	result := ""
	for _, d := range e.path {
		result = fmt.Sprintf("%s%s -> ", result, d)
	}
	if len(e.path) > 0 {
		result = fmt.Sprintf("%s%s", result, e.path[0])
	}
	return fmt.Sprintf("cyclic dependency: %s", result)
}

type InvalidConstructorError struct {
	cause string
}

func (e *InvalidConstructorError) Error() string {
	return fmt.Sprintf("invalid dependency constructor: %s", e.cause)
}

type DuplicatedNameError struct {
	name string
}

func (e *DuplicatedNameError) Error() string {
	return fmt.Sprintf("duplicated dependency name: %s", e.name)
}
