package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethodInt32 int32

const (
	PaymentMethodInt32Unspecified   PaymentMethodInt32 = 0
	PaymentMethodInt32Card          PaymentMethodInt32 = 1
	PaymentMethodInt32SBP           PaymentMethodInt32 = 2
	PaymentMethodInt32CreditCard    PaymentMethodInt32 = 3
	PaymentMethodInt32InvestorMoney PaymentMethodInt32 = 4
)

type PaymentMethodString string

const (
	PaymentMethodStringUnspecified   PaymentMethodString = "UNSPECIFIED"
	PaymentMethodStringCard          PaymentMethodString = "CARD"
	PaymentMethodStringSBP           PaymentMethodString = "SBP"
	PaymentMethodStringCreditCard    PaymentMethodString = "CREDIT_CARD"
	PaymentMethodStringInvestorMoney PaymentMethodString = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type GetOrderParams struct {
	OrderUUID uuid.UUID
}

type OrderDTO struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID
	WeaponUUID      *uuid.UUID
	TotalPrice      int64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethodString
	Status          OrderStatus
	CreatedAt       time.Time
}

type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID
	WeaponUUID      *uuid.UUID
	TotalPrice      int64
	TransactionUUID *uuid.UUID
	PaymentMethod   *PaymentMethodString
	Status          OrderStatus
	CreatedAt       time.Time
}

type CreateOrderRequest struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (r *CreateOrderRequest) PartUUIDs() []uuid.UUID {
	uuids := []uuid.UUID{r.HullUUID, r.EngineUUID}
	if r.ShieldUUID != nil {
		uuids = append(uuids, *r.ShieldUUID)
	}
	if r.WeaponUUID != nil {
		uuids = append(uuids, *r.WeaponUUID)
	}
	return uuids
}

type CreateOrderResponse struct {
	OrderUUID  uuid.UUID
	TotalPrice int64
}

type PayOrderRequest struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethodString
}

type PayOrderResponse struct {
	TransactionUUID uuid.UUID
}

type PayOrderClientRequest struct {
	OrderUUID     uuid.UUID
	PaymentMethod PaymentMethodInt32
}

type PayOrderClientResponse struct {
	TransactionUUID uuid.UUID
}
