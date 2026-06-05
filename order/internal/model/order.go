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
	OrderStatusAssembled      OrderStatus = "ASSEMBLED"
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

func (o *Order) SetStatus(status OrderStatus) {
	o.Status = status
}

func (o *Order) SetID(id string) {
	o.UUID = id
}

func (o *Order) SetTransactionID(transactionID string) {
	o.TransactionUUID = &transactionID
}

func (o *Order) SetPaymentMethod(method PaymentMethod) {
	o.PaymentMethod = &method
}

type OrderItem struct {
	UUID      string
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
