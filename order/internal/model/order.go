package model

import (
	"time"
)

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCreditCard    PaymentMethod = "CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

type Order struct {
	UUID            string
	Items           []OrderItem
	TotalPrice      int64
	TransactionUUID *string
	PaymentMethod   *PaymentMethod
	Status          OrderStatus
	CreatedAt       time.Time
}

type OrderItem struct {
	Uuid      string
	OrderUuid string
	PartUuid  string
	PartType  PartType
	Price     int64
	CreatedAt time.Time
}

type CreateOrderRequest struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID *string
	WeaponUUID *string
}

func (r *CreateOrderRequest) PartUUIDs() []string {
	uuids := []string{r.HullUUID, r.EngineUUID}
	if r.ShieldUUID != nil {
		uuids = append(uuids, *r.ShieldUUID)
	}
	if r.WeaponUUID != nil {
		uuids = append(uuids, *r.WeaponUUID)
	}
	return uuids
}

type PayOrderRequest struct {
	OrderUUID     string
	PaymentMethod PaymentMethod
}
