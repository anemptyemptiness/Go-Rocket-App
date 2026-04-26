package model

import "github.com/google/uuid"

type PaymentMethod int32

const (
	PaymentMethodUnspecified   PaymentMethod = 0
	PaymentMethodCard          PaymentMethod = 1
	PaymentMethodSBP           PaymentMethod = 2
	PaymentMethodCreditCard    PaymentMethod = 3
	PaymentMethodInvestorMoney PaymentMethod = 4
)

type PayOrderRequest struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}

type PayOrderResponse struct {
	TransactionUUID uuid.UUID
}
