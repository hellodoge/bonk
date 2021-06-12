package common

import (
	"errors"
	"fmt"
	"github.com/hellodoge/bonk/pkg/errors/wrapped"
	"runtime"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
)

func NewInvalidArgumentsError(inner error) wrapped.Error {
	// https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
	caller, _, _, ok := runtime.Caller(1)
	if !ok {
		return wrapped.Error{
			Outer: ErrInvalidArguments,
			Inner: inner,
		}
	}
	details := runtime.FuncForPC(caller)
	return wrapped.Error{
		Outer: fmt.Errorf("function %s: %w", details.Name(), ErrInvalidArguments),
		Inner: inner,
	}
}