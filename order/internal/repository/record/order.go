package record

import (
	"time"
)

type PaymentMethod string

type OrderStatus string

type Order struct {
	Uuid            string         `db:"uuid"`
	UserUuid        string         `db:"user_uuid"`
	TotalPrice      int64          `db:"total_price"`
	Status          OrderStatus    `db:"status"`
	TransactionUUID *string        `db:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `db:"payment_method"`
	OrderItems      []OrderItem    `json:"orderItems"`
	CreatedAt       time.Time      `db:"created_at"`
	UpdatedAt       *time.Time     `db:"updated_at"`
}

type PartType string

type OrderItem struct {
	Uuid      string    `db:"uuid" json:"uuid"`
	OrderUuid string    `db:"order_uuid" json:"orderUuid"`
	PartUuid  string    `db:"part_uuid" json:"partUuid"`
	PartType  PartType  `db:"part_type" json:"partType"`
	Price     int64     `db:"price" json:"price"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
