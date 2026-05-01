package model

import (
	"time"
)

type PartType string

func (pt *PartType) ToInt32() int32 {
	switch *pt {
	case PartTypeHull:
		return 1
	case PartTypeEngine:
		return 2
	case PartTypeShield:
		return 3
	case PartTypeWeapon:
		return 4
	default:
		return 0
	}
}

const (
	PartTypeUnspecified PartType = "UNSPECIFIED"
	PartTypeHull        PartType = "HULL"
	PartTypeEngine      PartType = "ENGINE"
	PartTypeShield      PartType = "SHIELD"
	PartTypeWeapon      PartType = "WEAPON"
)

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type ListPartsRequest struct {
	PartType PartType
	UUIDs    []string
}

type PartFilter struct {
	Uuids    []string
	PartType PartType
}
