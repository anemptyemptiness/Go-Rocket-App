package errors

import (
	"errors"
)

var (
	ErrPartNotFound                   = errors.New("деталь не найдена")
	ErrPartUUIDIsEmpty                = errors.New("идентификатор детали не может быть пустым")
	ErrEmptyRequest                   = errors.New("запрос не найден")
	ErrPartUUIDInvalid                = errors.New("идентификатор детали невалидный")
	ErrOutOfStock                     = errors.New("детали нет в наличии")
	ErrInvalidProperties              = errors.New("неверные свойства детали")
	ErrInvalidEngineClass             = errors.New("недопустимый класс двигателя")
	ErrInvalidShieldType              = errors.New("недопустимый тип щита")
	ErrInvalidWeaponType              = errors.New("недопустимый тип оружия")
	ErrNothingToRelease               = errors.New("нечего освобождать")
	ErrIncompatibleParts              = errors.New("детали несовместимы")
	ErrPartTypeMismatch               = errors.New("тип детали не соответствует слоту корабля")
	ErrStockQuantityOrReservedIsEmpty = errors.New("stock_quantity или reserved пустой")

	// AuthErrors.
	ErrEmptyMetadata    = errors.New("отсутствует metadata")
	ErrEmptySessionUUID = errors.New("отсутствует session-uuid")
	ErrExpiredSession   = errors.New("недействительная сессия")
)
