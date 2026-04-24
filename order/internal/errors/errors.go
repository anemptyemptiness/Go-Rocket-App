package errs

import (
	"errors"
	"fmt"
)

var (
	ErrOrderNotFound                    = errors.New("заказ не найден")
	ErrEmptyRequest                     = errors.New("запрос не должен быть пустым")
	ErrEmptyResponse                    = errors.New("ответ не должен быть пустым")
	ErrShieldUUIDIncorrect              = errors.New("shield_uuid некорректен")
	ErrWeaponUUIDIncorrect              = errors.New("weapon_uuid некорректен")
	ErrHullUUIDAndEngineUUIDAreRequired = errors.New("hull_uuid и engine_uuid обязательные параметры")
	ErrPartIsOver                       = errors.New("деталь закончилась")
	ErrUnknownPaymentMethod             = errors.New("неизвестный метод оплаты")
	ErrOrderStatusIncorrect             = errors.New("некорректный статус заказа")
	ErrOrderAlreadyPaid                 = errors.New("заказ уже оплачен")
	ErrOrderAlreadyCancelled            = errors.New("заказ уже отменён")
	ErrPaymentClientInvalidArgument     = errors.New("некорректный параметр")
	ErrPaymentClientInternal            = errors.New("внутренняя ошибка")
	ErrInventoryClientInternal          = errors.New("внутренняя ошибка")
)

func NewErrPartIsOver(partName string) error {
	return fmt.Errorf("%w: %s", ErrPartIsOver, partName)
}

func NewPaymentClientInvalidArgument(description string) error {
	return fmt.Errorf("%w: %s", ErrPaymentClientInvalidArgument, description)
}

func NewPaymentClientInternal(description string) error {
	return fmt.Errorf("%w: %s", ErrPaymentClientInternal, description)
}

func NewInventoryClientInternal(description string) error {
	return fmt.Errorf("%w: %s", ErrInventoryClientInternal, description)
}
