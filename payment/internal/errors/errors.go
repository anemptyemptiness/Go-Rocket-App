package errs

import "errors"

var (
	ErrEmptyRequest             = errors.New("запрос не найден")
	ErrOrderUUIDIsEmpty         = errors.New("идентификатор заказа не может быть пустым")
	ErrPaymentMethodUnspecified = errors.New("метод оплаты неопределённый")
)
