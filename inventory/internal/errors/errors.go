package errors

import (
	"errors"
)

var (
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrPartUUIDIsEmpty   = errors.New("идентификатор детали не может быть пустым")
	ErrEmptyRequest      = errors.New("запрос не найден")
	ErrIncorrectPartUUID = errors.New("невалидный UUID детали")
)
