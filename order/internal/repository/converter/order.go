package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func OrderRecordToModel(order record.Order) model.Order {
	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = new(model.PaymentMethod(*order.PaymentMethod))
	}

	orderItems := make([]model.OrderItem, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, model.OrderItem{
			UUID:      item.Uuid,
			OrderUuid: item.OrderUuid,
			PartUuid:  item.PartUuid,
			PartType:  model.PartType(item.PartType),
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return model.Order{
		UUID:            order.Uuid,
		UserUUID:        order.UserUuid,
		Items:           orderItems,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
