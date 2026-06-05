package errs

import (
	"fmt"
	"log/slog"
	"net/http"
)

// OgenError - constraint: ogen генерирует все error-типы как Error{Code, Message}.
type OgenError interface {
	~struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
}

type errorFactory func(code int, msg string) any

type httpErrorConfig struct {
	badRequest errorFactory
	notFound   errorFactory
	conflict   errorFactory
	internal   errorFactory
}

type Option func(*httpErrorConfig)

func WithBadRequest[T OgenError]() Option {
	return func(c *httpErrorConfig) { c.badRequest = newFactory[T]() }
}

func WithNotFound[T OgenError]() Option {
	return func(c *httpErrorConfig) { c.notFound = newFactory[T]() }
}

func WithConflict[T OgenError]() Option {
	return func(c *httpErrorConfig) { c.conflict = newFactory[T]() }
}

func WithInternal[T OgenError]() Option {
	return func(c *httpErrorConfig) { c.internal = newFactory[T]() }
}

func newFactory[T OgenError]() errorFactory {
	return func(code int, msg string) any {
		r := T(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{Code: code, Message: msg})
		return &r
	}
}

// MapHTTPError конвертирует BusinessError в типизированный ogen-ответ.
// TRes - интерфейс ogen-ответа (например orderv1.PayOrderRes).
func MapHTTPError[TRes any](err error, opts ...Option) (TRes, error) {
	cfg := &httpErrorConfig{}
	for _, o := range opts {
		o(cfg)
	}

	result := buildResponse(err, cfg)

	res, ok := result.(TRes)
	if !ok {
		var zero TRes
		return zero, fmt.Errorf("unexpected error response type: %T", result)
	}
	return res, nil
}

func buildResponse(err error, cfg *httpErrorConfig) any {
	be, ok := AsBusinessError(err)
	if !ok {
		slog.Error("unexpected error", "error", err)
		return cfg.fallback(http.StatusInternalServerError, "внутренняя ошибка")
	}

	code := be.Code().HTTPCode()
	msg := be.ClientMessage()

	switch be.Code() {
	case CodeInvalidArgument:
		if cfg.badRequest != nil {
			return cfg.badRequest(code, msg)
		}
	case CodeNotFound:
		if cfg.notFound != nil {
			return cfg.notFound(code, msg)
		}
	case CodeConflict:
		if cfg.conflict != nil {
			return cfg.conflict(code, msg)
		}
	}

	return cfg.fallback(code, msg)
}

func (c *httpErrorConfig) fallback(code int, msg string) any {
	if c.internal != nil {
		return c.internal(code, msg)
	}
	slog.Error("MapHTTPError: no internal error factory configured")
	return nil
}
