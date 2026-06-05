package valueobject

import (
	"fmt"
	"strings"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

type PartType string

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

func NewPartType(s string) (PartType, error) {
	pt := PartType(s)
	switch pt {
	case PartTypeHull, PartTypeEngine, PartTypeShield, PartTypeWeapon, PartTypeUnspecified:
		return pt, nil
	default:
		return "", pkgerr.Internal(fmt.Errorf("неизвестный тип детали %q: %w", s, errs.ErrInvalidProperties))
	}
}

func (pt PartType) ToProto() inventoryv1.PartType {
	return inventoryv1.PartType(inventoryv1.PartType_value["PART_TYPE_"+strings.ToUpper(string(pt))])
}
