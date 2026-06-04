package valueobject

import (
	"fmt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

type HullProperties interface {
	Strength() int32
	CanSupport(e EngineProperties) bool
}

type hullProperties struct {
	strength int32
}

func (hp *hullProperties) Strength() int32 {
	return hp.strength
}

func (hp *hullProperties) CanSupport(e EngineProperties) bool {
	return hp.Strength() >= e.RequiredStrength()
}

// NewHullProperties создаёт свойства корпуса. Прочность должна быть в диапазоне 30–200.
func NewHullProperties(strength int32) (PartProperties, error) {
	if strength < 30 || strength > 200 {
		return nil, pkgerr.Internal(fmt.Errorf("прочность корпуса должна быть от 30 до 200, получено %d: %w", strength, errs.ErrInvalidProperties))
	}
	return &partProperties{
		hull: &hullProperties{strength: strength},
	}, nil
}
