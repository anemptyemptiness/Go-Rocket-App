package compatibility_checker

import (
	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

type partSet struct {
	hull   valueobject.HullProperties
	engine valueobject.EngineProperties
	shield valueobject.ShieldProperties
	weapon valueobject.WeaponProperties
}

func (s *service) Check(parts []entity.Part) error {
	set := extractProperties(parts)

	if err := checkHullEngine(set); err != nil {
		return err
	}

	if err := checkShieldWeapon(set); err != nil {
		return err
	}

	return nil
}

func extractProperties(parts []entity.Part) partSet {
	var set partSet

	for _, part := range parts {
		props := part.GetProperties()
		if props != nil && props.Hull() != nil {
			set.hull = props.Hull()
		}
		if props != nil && props.Engine() != nil {
			set.engine = props.Engine()
		}
		if props != nil && props.Shield() != nil {
			set.shield = props.Shield()
		}
		if props != nil && props.Weapon() != nil {
			set.weapon = props.Weapon()
		}
	}

	return set
}

func checkHullEngine(set partSet) error {
	if set.engine.Class() == valueobject.EngineClassC &&
		!set.hull.CanSupport(set.engine) &&
		set.hull.Strength() < 30 {
		return pkgerr.FailedPrecondition(errs.ErrIncompatibleParts)
	}

	if set.engine.Class() == valueobject.EngineClassB &&
		!set.hull.CanSupport(set.engine) &&
		set.hull.Strength() < 70 {
		return pkgerr.FailedPrecondition(errs.ErrIncompatibleParts)
	}

	if set.engine.Class() == valueobject.EngineClassA &&
		!set.hull.CanSupport(set.engine) &&
		set.hull.Strength() < 100 {
		return pkgerr.FailedPrecondition(errs.ErrIncompatibleParts)
	}

	return nil
}

func checkShieldWeapon(set partSet) error {
	if set.shield != nil && set.weapon != nil && set.shield.ConflictsWith(set.weapon) {
		return pkgerr.FailedPrecondition(errs.ErrIncompatibleParts)
	}

	return nil
}
