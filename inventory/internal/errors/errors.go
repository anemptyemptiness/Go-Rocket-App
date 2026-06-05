package errors

import (
	"errors"
)

var (
	ErrPartNotFound    = errors.New("деталь не найдена")
	ErrPartUUIDIsEmpty = errors.New("идентификатор детали не может быть пустым")
	ErrEmptyRequest    = errors.New("запрос не найден")
	ErrPartUUIDInvalid = errors.New("идентификатор детали невалидный")
)
