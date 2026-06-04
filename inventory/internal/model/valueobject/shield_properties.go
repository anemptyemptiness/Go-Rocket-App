package valueobject

import (
	"fmt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
)

type ShieldType string

const (
	ShieldTypeEnergy ShieldType = "energy"
	ShieldTypePlasma ShieldType = "plasma"
)

func NewShieldType(s string) (ShieldType, error) {
	st := ShieldType(s)

	if st != ShieldTypeEnergy && st != ShieldTypePlasma {
		return "", fmt.Errorf("неверный тип щита (%s): %w", st, errs.ErrInvalidShieldType)
	}
	return st, nil
}

type ShieldProperties interface {
	Type() ShieldType
	ConflictsWith(w WeaponProperties) bool
}

type shieldProperties struct {
	shieldType ShieldType
}

func NewShieldProperties(shieldType ShieldType) (PartProperties, error) {
	return &partProperties{
		shield: &shieldProperties{shieldType: shieldType},
	}, nil
}

func (sp *shieldProperties) Type() ShieldType {
	return sp.shieldType
}

func (sp *shieldProperties) ConflictsWith(w WeaponProperties) bool {
	return sp.Type() == ShieldTypePlasma && w.Type() == WeaponTypeLaser
}
