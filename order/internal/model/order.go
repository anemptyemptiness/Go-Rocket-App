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

type CreateOrderResponse struct {
	OrderUUID  uuid.UUID
	TotalPrice int64
}

type ListPartsClientRequest struct {
	UUIDs []uuid.UUID
}

type PartType int32

const (
	PartTypeUnspecified PartType = 0
	PartTypeHull        PartType = 1
	PartTypeEngine      PartType = 2
	PartTypeShield      PartType = 3
	PartTypeWeapon      PartType = 4
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type ListPartsClientResponse struct {
	Parts []Part
}

type PayOrderRequest struct {
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
