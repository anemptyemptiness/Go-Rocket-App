package errors

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrEmptyRequest      = errors.New("запрос не найден")
	ErrIncorrectPartUUID = errors.New("невалидный UUID детали")
)

func NewErrPartNotFound(partUUID uuid.UUID) error {
	return fmt.Errorf("%w: %s", ErrPartNotFound, partUUID.String())
}
