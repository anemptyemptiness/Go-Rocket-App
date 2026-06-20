package converter

import (
	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/input"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
)

func OrderModelToDTO(order model.Order) *orderv1.OrderDto {
	var hullUuid, engineUuid uuid.UUID
	var shieldUuid, weaponUuid orderv1.OptNilUUID

	for _, item := range order.Items {
		switch item.PartType {
		case model.PartTypeHull:
			hullUuid = uuid.MustParse(item.PartUuid)
		case model.PartTypeEngine:
			engineUuid = uuid.MustParse(item.PartUuid)
		case model.PartTypeShield:
			shieldUuid = orderv1.NewOptNilUUID(uuid.MustParse(item.PartUuid))
		case model.PartTypeWeapon:
			weaponUuid = orderv1.NewOptNilUUID(uuid.MustParse(item.PartUuid))
		}
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(uuid.MustParse(*order.TransactionUUID))
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       uuid.MustParse(order.UUID),
		UserUUID:        uuid.MustParse(order.UserUUID),
		HullUUID:        hullUuid,
		EngineUUID:      engineUuid,
		ShieldUUID:      shieldUuid,
		WeaponUUID:      weaponUuid,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

func CreateOrderRequestToModel(req *orderv1.CreateOrderRequest) (input.CreateOrderRequest, error) {
	if req == nil {
		return input.CreateOrderRequest{}, pkgerr.InvalidArgument(errs.ErrEmptyRequest)
	}

	if req.GetHullUUID().String() == "" || req.GetHullUUID() == uuid.Nil ||
		req.GetEngineUUID().String() == "" || req.GetEngineUUID() == uuid.Nil {
		return input.CreateOrderRequest{}, pkgerr.InvalidArgument(errs.ErrHullUUIDAndEngineUUIDAreRequired)
	}

	var shieldUUID *string
	if v := req.GetShieldUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return input.CreateOrderRequest{}, pkgerr.InvalidArgument(errs.ErrShieldUUIDIncorrect)
		}

		shieldUUID = new(id.String())
	}

	var weaponUUID *string
	if v := req.GetWeaponUUID(); v.IsSet() && !v.IsNull() {
		id, ok := v.Get()
		if !ok {
			return input.CreateOrderRequest{}, pkgerr.InvalidArgument(errs.ErrWeaponUUIDIncorrect)
		}

		weaponUUID = new(id.String())
	}

	return input.CreateOrderRequest{
		HullUUID:   req.GetHullUUID().String(),
		EngineUUID: req.GetEngineUUID().String(),
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
	}, nil
}

func PayOrderRequestToModel(req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (input.PayOrderRequest, error) {
	if req == nil {
		return input.PayOrderRequest{}, pkgerr.InvalidArgument(errs.ErrEmptyRequest)
	}

	return input.PayOrderRequest{
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
		OrderUUID:     params.OrderUUID.String(),
	}, nil
}
