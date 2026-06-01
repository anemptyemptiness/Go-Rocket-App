package valueobject

import (
	"fmt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

type EngineClass string

const (
	EngineClassA EngineClass = "A"
	EngineClassB EngineClass = "B"
	EngineClassC EngineClass = "C"
)

func NewEngineClass(s string) (EngineClass, error) {
	ec := EngineClass(s)
	switch ec {
	case EngineClassA, EngineClassB, EngineClassC:
		return ec, nil
	default:
		return "", fmt.Errorf("недопустимый класс двигателя: %w", errs.ErrInvalidEngineClass)
	}
}

type EngineProperties interface {
	Class() EngineClass
	RequiredStrength() int32
}

type engineProperties struct {
	class            EngineClass
	requiredStrength int32
}

func (ep *engineProperties) Class() EngineClass {
	return ep.class
}

func (ep *engineProperties) RequiredStrength() int32 {
	return ep.requiredStrength
}

func NewEngineProperties(class EngineClass, requiredStrength int32) (PartProperties, error) {
	if requiredStrength <= 0 {
		return nil, pkgerr.Internal(fmt.Errorf("требуемая сила двигателя должна быть положительной: %w", errs.ErrInvalidProperties))
	}
	if class == EngineClassC && requiredStrength != 30 {
		return nil, pkgerr.Internal(fmt.Errorf("для класса C требуемая сила должна быть 30: %w", errs.ErrInvalidProperties))
	}
	if class == EngineClassB && requiredStrength != 70 {
		return nil, pkgerr.Internal(fmt.Errorf("для класса B требуемая сила должна быть 70: %w", errs.ErrInvalidProperties))
	}
	return &partProperties{
		engine: &engineProperties{
			class:            class,
			requiredStrength: requiredStrength,
		},
	}, nil
}
