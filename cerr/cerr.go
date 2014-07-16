package cerr

import (
	"fmt"
)

type BadRequestError struct {
	s string
}

func NewBadRequestError(format string, a ...interface{}) error {
	return &BadRequestError{s: fmt.Sprintf(format, a...)}
}

func (e *BadRequestError) Error() string {
	return e.s
}

type NotFoundError struct {
	s string
}

func NewNotFoundError(format string, a ...interface{}) error {
	return &NotFoundError{s: fmt.Sprintf(format, a...)}
}

func (e *NotFoundError) Error() string {
	return e.s
}

type ServerError struct {
	s string
}

func NewServerError(format string, a ...interface{}) error {
	return &ServerError{s: fmt.Sprintf(format, a...)}
}

func (e *ServerError) Error() string {
	return e.s
}
