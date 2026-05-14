package part

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/record"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	const query = `
		SELECT 
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.created_at
		FROM parts AS p
		WHERE p.uuid = $1;`

	var part record.Part

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, query, uuid).Scan(
		&part.UUID,
		&part.Name,
		&part.Description,
		&part.PartType,
		&part.Price,
		&part.StockQuantity,
		&part.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, errs.ErrPartNotFound) {
			return model.Part{}, fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, uuid)
		}
		return model.Part{}, err
	}

	return repoConverter.PartRecordToModel(part), nil
}

func (r *repository) ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	query := `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.created_at
		FROM parts AS p `

	switch {
	case len(filter.UUIDs) > 0:
		query += `WHERE p.uuid = ANY($1::UUID[]);`
	case filter.PartType != "":
		query += `
			WHERE p.part_type = $2
			ORDER BY p.name;`
	default:
		query += `;`
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, filter.UUIDs, filter.PartType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: filter=%v", errs.ErrPartNotFound, filter)
		}
		return nil, err
	}
	defer rows.Close()

	parts, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.Part])
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, errs.ErrPartNotFound
	}

	return repoConverter.PartsRecordToModel(parts), nil
}
