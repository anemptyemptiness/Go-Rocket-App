package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func OrderModelToDTO(order model.Order) *orderv1.OrderDto {
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func CreateOrderRequestToModel(req *orderv1.CreateOrderRequest) (model.CreateOrderRequest, error) {
	if req == nil {
		return model.CreateOrderRequest{}, errs.ErrEmptyRequest
	}

	if req.GetHullUUID() == uuid.Nil || req.GetEngineUUID() == uuid.Nil {
		return model.CreateOrderRequest{}, errs.ErrHullUUIDAndEngineUUIDAreRequired
	}

	var shieldUUID *uuid.UUID
	if v := req.GetShieldUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return model.CreateOrderRequest{}, errs.ErrShieldUUIDIncorrect
		}
		shieldUUID = &id
	}

	var weaponUUID *uuid.UUID
	if v := req.GetWeaponUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return model.CreateOrderRequest{}, errs.ErrWeaponUUIDIncorrect
		}
		weaponUUID = &id
	}

	return model.CreateOrderRequest{
		HullUUID:   req.GetHullUUID(),
		EngineUUID: req.GetEngineUUID(),
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}, nil
}

func PayOrderRequestToModel(req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (model.PayOrderRequest, error) {
	if req == nil {
		return model.PayOrderRequest{}, errs.ErrEmptyRequest
	}

	return model.PayOrderRequest{
		PaymentMethod: model.PaymentMethodString(req.GetPaymentMethod()),
		OrderUUID:     params.OrderUUID,
	}, nil
}

func PaymentMethodFromInt32ToString(paymentMethod model.PaymentMethodInt32) model.PaymentMethodString {
	switch paymentMethod {
	case model.PaymentMethodInt32Card:
		return model.PaymentMethodStringCard
	case model.PaymentMethodInt32SBP:
		return model.PaymentMethodStringSBP
	case model.PaymentMethodInt32CreditCard:
		return model.PaymentMethodStringCreditCard
	default:
		return model.PaymentMethodStringInvestorMoney
	}
}
