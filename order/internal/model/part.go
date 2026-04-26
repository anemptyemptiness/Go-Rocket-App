package model

import (
	"time"

	"github.com/google/uuid"
)

type ListPartsClientRequest struct {
	UUIDs []uuid.UUID
}

type PartType int32

const (
	PartTypeUnspecified PartType = 0
	PartTypeHull        PartType = 1
	PartTypeEngine      PartType = 2
	PartTypeShield      PartType = 3
	PartTypeWeapon      PartType = 4
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      PartType
	StockQuantity int64
	CreatedAt     time.Time
}

type ListPartsClientResponse struct {
	Parts []Part
}
