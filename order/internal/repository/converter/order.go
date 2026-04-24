package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func OrderRecordToModel(order record.Order) model.Order {
	paymentMethod := model.PaymentMethodStringUnspecified
	if order.PaymentMethod != nil {
		paymentMethod = model.PaymentMethodString(*order.PaymentMethod)
	}

	return model.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func OrderModelToRecord(order model.Order) record.Order {
	paymentMethod := record.PaymentMethodUnspecified
	if order.PaymentMethod != nil {
		paymentMethod = record.PaymentMethod(*order.PaymentMethod)
	}

	return record.Order{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      order.ShieldUUID,
		WeaponUUID:      order.WeaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          record.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
