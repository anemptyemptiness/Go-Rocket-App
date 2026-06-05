package entity

import (
	"time"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

type Part struct {
	partUUID      string
	name          string
	description   string
	price         int64
	stockQuantity int64
	reserved      int64
	partType      valueobject.PartType
	properties    valueobject.PartProperties
	createdAt     time.Time
}

func (p *Part) GetPartUUID() string {
	return p.partUUID
}

func (p *Part) GetName() string {
	return p.name
}

func (p *Part) GetDescription() string {
	return p.description
}

func (p *Part) GetPrice() int64 {
	return p.price
}

func (p *Part) GetPartType() valueobject.PartType {
	return p.partType
}

func (p *Part) GetStockQuantity() int64 {
	return p.stockQuantity
}

func (p *Part) GetCreatedAt() time.Time {
	return p.createdAt
}

func (p *Part) GetProperties() valueobject.PartProperties {
	return p.properties
}

func (p *Part) GetReserved() int64 {
	return p.reserved
}

func (p *Part) Reserve() error {
	if p.Available() <= 0 {
		return pkgerr.ResourceExhausted(errs.ErrOutOfStock)
	}

	p.reserved++

	return nil
}

func (p *Part) Available() int64 {
	return p.stockQuantity - p.reserved
}

func (p *Part) Release() error {
	if p.reserved <= 0 {
		return pkgerr.FailedPrecondition(errs.ErrNothingToRelease)
	}

	p.reserved--

	return nil
}

func RestorePart(
	uuid, name, description string,
	price, stockQuantity, reserved int64,
	partType valueobject.PartType,
	properties valueobject.PartProperties,
	createdAt time.Time,
) Part {
	return Part{
		partUUID:      uuid,
		name:          name,
		description:   description,
		price:         price,
		stockQuantity: stockQuantity,
		reserved:      reserved,
		partType:      partType,
		properties:    properties,
		createdAt:     createdAt,
	}
}
