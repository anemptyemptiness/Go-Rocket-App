package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func OrderRecordToModel(order record.Order) model.Order {
	var paymentMethod *model.PaymentMethodString
	if order.PaymentMethod != nil {
		value := model.PaymentMethodString(*order.PaymentMethod)
		paymentMethod = &value
	}

	return model.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderModelToRecord(order model.Order) record.Order {
	var paymentMethod *record.PaymentMethod
	if order.PaymentMethod != nil {
		value := record.PaymentMethod(*order.PaymentMethod)
		paymentMethod = &value
	}

	return record.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          record.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
