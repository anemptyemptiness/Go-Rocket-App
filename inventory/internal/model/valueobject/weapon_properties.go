package valueobject

import (
	"fmt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
)

type Type string

const (
	WeaponTypeLaser   Type = "laser"
	WeaponTypeMissile Type = "missile"
)

func NewWeaponType(s string) (Type, error) {
	wt := Type(s)

	if wt != WeaponTypeLaser && wt != WeaponTypeMissile {
		return "", fmt.Errorf("неверный тип оружия (%s): %w", wt, errs.ErrInvalidWeaponType)
	}
	return wt, nil
}

type WeaponProperties interface {
	Type() Type
}

type weaponProperties struct {
	weaponType Type
}

func NewWeaponProperties(weaponType Type) (PartProperties, error) {
	return &partProperties{
		weapon: &weaponProperties{weaponType: weaponType},
	}, nil
}

func (wp *weaponProperties) Type() Type {
	return wp.weaponType
}
