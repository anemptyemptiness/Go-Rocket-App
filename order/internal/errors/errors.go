package errs

import (
	"errors"
)

var (
	// OrderErrors.
	ErrOrderNotFound                    = errors.New("заказ не найден")
	ErrEmptyRequest                     = errors.New("запрос не должен быть пустым")
	ErrEmptyResponse                    = errors.New("ответ не должен быть пустым")
	ErrShieldUUIDIncorrect              = errors.New("shield_uuid некорректен")
	ErrWeaponUUIDIncorrect              = errors.New("weapon_uuid некорректен")
	ErrHullUUIDAndEngineUUIDAreRequired = errors.New("hull_uuid и engine_uuid обязательные параметры")
	ErrUserUUIDIsRequired               = errors.New("user_uuid обязательный параметр")
	ErrPartIsOver                       = errors.New("деталь закончилась")
	ErrOrderStatusIncorrect             = errors.New("некорректный статус заказа")
	ErrOrderAlreadyPaid                 = errors.New("заказ уже оплачен")
	ErrOrderAlreadyCancelled            = errors.New("заказ уже отменён")
	ErrOrderAssembled                   = errors.New("заказ уже в сборке")
	ErrInvalidPartUUID                  = errors.New("идентификатор детали невалидный")
)
