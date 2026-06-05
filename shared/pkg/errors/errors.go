package errs

import (
	"errors"
	"net/http"

	"google.golang.org/grpc/codes"
)

type Code string

const (
	CodeUnknown         Code = "unknown"
	CodeInvalidArgument Code = "invalid_argument"
	CodeNotFound        Code = "not_found"
	CodeConflict        Code = "conflict"
	CodeForbidden       Code = "forbidden"
	CodeUnauthorized    Code = "unauthorized"
	CodeInternal        Code = "internal"
)

var (
	ErrInternal     = errors.New("internal error")
	ErrInvalidArg   = errors.New("invalid argument")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrForbidden    = errors.New("forbidden")
	ErrUnauthorized = errors.New("unauthorized")
)

type BusinessError struct {
	code    Code
	cause   error
	message string
}

func (e *BusinessError) Error() string {
	if e == nil {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return string(e.code)
}

func (e *BusinessError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *BusinessError) Code() Code {
	if e == nil {
		return CodeUnknown
	}
	return e.code
}

func (e *BusinessError) ClientMessage() string {
	if e == nil {
		return ""
	}
	if e.message != "" {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return string(e.code)
}

func (e *BusinessError) Is(target error) bool {
	var t *BusinessError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.code == t.code
}

func NewBusinessError(code Code, err error) error {
	if err == nil {
		return nil
	}
	return &BusinessError{code: code, cause: err}
}

func NewBusinessErrorMessage(code Code, err error, message string) error {
	if err == nil {
		return nil
	}
	return &BusinessError{code: code, cause: err, message: message}
}

func InvalidArgument(err error) error { return NewBusinessError(CodeInvalidArgument, err) }
func NotFound(err error) error        { return NewBusinessError(CodeNotFound, err) }
func Conflict(err error) error        { return NewBusinessError(CodeConflict, err) }
func Forbidden(err error) error       { return NewBusinessError(CodeForbidden, err) }
func Unauthorized(err error) error    { return NewBusinessError(CodeUnauthorized, err) }
func Internal(err error) error        { return NewBusinessError(CodeInternal, err) }

func IsBusinessError(err error) bool {
	var be *BusinessError
	return errors.As(err, &be)
}

func AsBusinessError(err error) (*BusinessError, bool) {
	var be *BusinessError
	if !errors.As(err, &be) {
		return nil, false
	}
	return be, true
}

func (c Code) GRPCCode() codes.Code {
	switch c {
	case CodeInvalidArgument:
		return codes.InvalidArgument
	case CodeNotFound:
		return codes.NotFound
	case CodeConflict:
		return codes.AlreadyExists
	case CodeForbidden:
		return codes.PermissionDenied
	case CodeUnauthorized:
		return codes.Unauthenticated
	case CodeInternal:
		return codes.Internal
	default:
		return codes.Internal
	}
}

func (c Code) HTTPCode() int {
	switch c {
	case CodeInvalidArgument:
		return http.StatusBadRequest
	case CodeNotFound:
		return http.StatusNotFound
	case CodeConflict:
		return http.StatusConflict
	case CodeForbidden:
		return http.StatusForbidden
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
