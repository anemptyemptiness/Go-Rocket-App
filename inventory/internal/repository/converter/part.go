package converter

import (
	"encoding/json"
	"fmt"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/record"
)

func PartsRecordToModel(parts []record.Part) ([]entity.Part, error) {
	modelParts := make([]entity.Part, 0, len(parts))

	for _, part := range parts {
		modelPart, err := PartRecordToModel(part)
		if err != nil {
			return nil, err
		}

		modelParts = append(modelParts, modelPart)
	}

	return modelParts, nil
}

func PartRecordToModel(part record.Part) (entity.Part, error) {
	pt, err := valueobject.NewPartType(part.PartType)
	if err != nil {
		return entity.Part{}, fmt.Errorf("конвертировать тип детали: %w", err)
	}

	var partPropertiesRecord record.PartProperties
	if err = json.Unmarshal(part.PartProperties, &partPropertiesRecord); err != nil {
		return entity.Part{}, fmt.Errorf("десериализовать свойства детали: %w", err)
	}

	props, err := partPropertiesFromRecord(partPropertiesRecord)
	if err != nil {
		return entity.Part{}, fmt.Errorf("конвертировать свойства детали: %w", err)
	}

	return entity.RestorePart(
		part.UUID,
		part.Name,
		part.Description,
		part.Price,
		part.StockQuantity,
		part.Reserved,
		pt,
		props,
		part.CreatedAt,
	), nil
}

func partPropertiesFromRecord(rec record.PartProperties) (valueobject.PartProperties, error) {
	switch {
	case rec.Engine != nil:
		engineClass, err := valueobject.NewEngineClass(rec.Engine.Class)
		if err != nil {
			return nil, fmt.Errorf("конвертировать класс двигателя: %w", err)
		}

		return valueobject.NewEngineProperties(
			engineClass,
			rec.Engine.RequiredStrength,
		)
	case rec.Hull != nil:
		return valueobject.NewHullProperties(
			rec.Hull.Strength,
		)
	case rec.Shield != nil:
		shieldType, err := valueobject.NewShieldType(rec.Shield.ShieldType)
		if err != nil {
			return nil, fmt.Errorf("конвертировать тип щита: %w", err)
		}

		return valueobject.NewShieldProperties(
			shieldType,
		)
	case rec.Weapon != nil:
		weaponType, err := valueobject.NewWeaponType(rec.Weapon.WeaponType)
		if err != nil {
			return nil, fmt.Errorf("конвертировать тип оружия: %w", err)
		}

		return valueobject.NewWeaponProperties(
			weaponType,
		)
	default:
		return nil, errs.ErrInvalidProperties
	}
}
