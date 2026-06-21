package input

import "github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"

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
	PaymentMethod model.PaymentMethod
}
