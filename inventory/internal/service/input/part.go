package input

import "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"

type PartFilter struct {
	UUIDs    []string
	PartType valueobject.PartType
}

type ValidateCompatibilityRequest struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID string
	WeaponUUID string
}

func (v ValidateCompatibilityRequest) UUIDs() []string {
	uuids := []string{v.HullUUID, v.EngineUUID}

	if v.ShieldUUID != "" {
		uuids = append(uuids, v.ShieldUUID)
	}
	if v.WeaponUUID != "" {
		uuids = append(uuids, v.WeaponUUID)
	}

	return uuids
}

type ReservePartsRequest struct {
	UUIDs []string
}

type ReleasePartsRequest struct {
	UUIDs []string
}
