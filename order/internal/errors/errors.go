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
	ErrPartIsOver                       = errors.New("деталь закончилась")
	ErrOrderStatusIncorrect             = errors.New("некорректный статус заказа")
	ErrOrderAlreadyPaid                 = errors.New("заказ уже оплачен")
	ErrOrderAlreadyCancelled            = errors.New("заказ уже отменён")

	// Payment Client Errors.
	ErrPaymentClientInvalidArgument = errors.New("некорректный(ые) параметр(ы)")

	// Inventory Client Errors.
	ErrInventoryClientInvalidArgument = errors.New("некорректный(ые) параметр(ы)")
	ErrInventoryClientNotFound        = errors.New("элемент(ы) не найден(ы)")
)
