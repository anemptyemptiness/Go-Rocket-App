package converter

import (
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func OrderRecordToModel(order record.Order) model.Order {
	var paymentMethod *model.PaymentMethod
	if order.PaymentMethod != nil {
		value := model.PaymentMethod(*order.PaymentMethod)
		paymentMethod = &value
	}

	orderItems := make([]model.OrderItem, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		orderItems = append(orderItems, model.OrderItem{
			Uuid:      item.Uuid,
			OrderUuid: item.OrderUuid,
			PartUuid:  item.PartUuid,
			PartType:  model.PartType(item.PartType),
			Price:     item.Price,
			CreatedAt: item.CreatedAt,
		})
	}

	return model.Order{
		UUID:            order.Uuid,
		Items:           orderItems,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}
