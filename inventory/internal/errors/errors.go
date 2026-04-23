package errors

import "errors"

var (
	ErrPartNotFound = errors.New("деталь не найдена")
	ErrEmptyRequest = errors.New("запрос не найден")
)
