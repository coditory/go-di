package di

import (
	"errors"
	"fmt"
	"reflect"
)

var ErrSkippedDependency = errors.New("skipped dependency")

const (
	ErrTypeDependencyCreation = iota
	ErrTypeDuplicatedName
	ErrTypeDuplicatedRegistration
	ErrTypeMissingDependency
	ErrTypeInvalidType
	ErrTypeInvalidConstructor
	ErrTypeCyclicDependency
	ErrTypeDependencyInitialization
	ErrTypeDependencyShutdown
	ErrTypeLifecycle
)

type Error struct {
	errType int
	message string
	cause   error
}

func (e *Error) ErrType() int {
	return e.errType
}

func (e *Error) IsErrType(errType int) bool {
	return e.errType == errType
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) RootCause() error {
	if e.cause == nil {
		return e
	}
	if dierr, ok := e.cause.(*Error); ok {
		return dierr.RootCause()
	}
	return e.cause
}

func newLifecycleError(cause string) *Error {
	msg := fmt.Sprintf("context lifecycle error: %s", cause)
	return &Error{
		errType: ErrTypeLifecycle,
		message: msg,
	}
}

func newDuplicatedRegistrationError() *Error {
	return &Error{
		errType: ErrTypeDuplicatedRegistration,
		message: "duplicated registration",
	}
}

func newDependencyCreationError(objName *string, objType *reflect.Type, cause error) *Error {
	msg := fmt.Sprintf("could not create dependency %s, cause:\n%s", descriptor(objName, objType), cause)
	return &Error{
		errType: ErrTypeDependencyCreation,
		message: msg,
		cause:   cause,
	}
}

func newMissingDependencyError(objName *string, objType *reflect.Type) *Error {
	msg := fmt.Sprintf("missing dependency %s", descriptor(objName, objType))
	return &Error{
		errType: ErrTypeMissingDependency,
		message: msg,
	}
}

func newInvalidTypeError(objName *string, objType reflect.Type, expectedType reflect.Type) *Error {
	msg := fmt.Sprintf("could not cast %s to %s", descriptor(objName, &objType), expectedType)
	return &Error{
		errType: ErrTypeInvalidType,
		message: msg,
	}
}

func newInitializationError(objType *reflect.Type, cause error) *Error {
	msg := fmt.Sprintf("could not initialize dependency: %s, cause:\n%s", descriptor(nil, objType), cause)
	return &Error{
		errType: ErrTypeDependencyInitialization,
		message: msg,
		cause:   cause,
	}
}

func newShutdownError(objType *reflect.Type, cause error) *Error {
	msg := fmt.Sprintf("could not shutdown dependency: %s, cause:\n%s", descriptor(nil, objType), cause)
	return &Error{
		errType: ErrTypeDependencyShutdown,
		message: msg,
		cause:   cause,
	}
}

func newCyclicDependencyError(path []string) *Error {
	msg := ""
	for _, d := range path {
		msg = fmt.Sprintf("%s%s -> ", msg, d)
	}
	if len(path) > 0 {
		msg = fmt.Sprintf("%s%s", msg, path[0])
	}
	msg = fmt.Sprintf("cyclic dependency: %s", msg)
	return &Error{
		errType: ErrTypeCyclicDependency,
		message: msg,
	}
}

func newInvalidConstructorError(cause string) *Error {
	msg := fmt.Sprintf("invalid dependency constructor: %s", cause)
	return &Error{
		errType: ErrTypeInvalidConstructor,
		message: msg,
	}
}

func newDuplicatedNameError(name string) *Error {
	msg := fmt.Sprintf("duplicated dependency name: %s", name)
	return &Error{
		errType: ErrTypeDuplicatedName,
		message: msg,
	}
}
