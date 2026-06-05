package model

import (
	"strings"
	"time"

	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
)

type PartType string

func (pt PartType) ToProto() inventoryv1.PartType {
	return inventoryv1.PartType(inventoryv1.PartType_value["PART_TYPE_"+strings.ToUpper(string(pt))])
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

type PartFilter struct {
	UUIDs    []string
	PartType PartType
}
