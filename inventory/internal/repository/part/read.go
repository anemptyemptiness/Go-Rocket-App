package part

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/valueobject"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/repository/record"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (entity.Part, error) {
	const query = `
		SELECT 
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.reserved,
			p.properties,
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
		&part.Reserved,
		&part.PartProperties,
		&part.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Part{}, fmt.Errorf("%w: uuid=%s", errs.ErrPartNotFound, uuid)
		}
		return entity.Part{}, err
	}

	return repoConverter.PartRecordToModel(part)
}

//nolint:dupl // Два дублирующих метода созданы намеренно, так как ListForUpdate используется для транзакционности операции.
func (r *repository) ListParts(ctx context.Context, filter input.PartFilter) ([]entity.Part, error) {
	query := `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.reserved,
			p.properties,
			p.created_at
		FROM parts AS p `

	var args []any

	switch {
	case len(filter.UUIDs) > 0:
		query += `
			WHERE p.uuid = ANY($1::UUID[])
			ORDER BY ARRAY_POSITION($1::UUID[], p.uuid);`
		args = append(args, filter.UUIDs)
	case filter.PartType != "" && filter.PartType != valueobject.PartTypeUnspecified:
		query += `
			WHERE p.part_type = $1::VARCHAR
			ORDER BY p.name ASC;`
		args = append(args, string(filter.PartType))
	default:
		query += `ORDER BY p.name ASC;`
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, args...)
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

	return repoConverter.PartsRecordToModel(parts)
}

//nolint:dupl // Два дублирующих метода созданы намеренно, так как ListForUpdate используется для транзакционности операции.
func (r *repository) ListForUpdate(ctx context.Context, filter input.PartFilter) ([]entity.Part, error) {
	query := `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.reserved,
			p.properties,
			p.created_at
		FROM parts AS p `

	var args []any

	switch {
	case len(filter.UUIDs) > 0:
		query += `
			WHERE p.uuid = ANY($1::UUID[])
			ORDER BY p.uuid FOR UPDATE;`
		args = append(args, filter.UUIDs)
	case filter.PartType != "" && filter.PartType != valueobject.PartTypeUnspecified:
		query += `
			WHERE p.part_type = $1::VARCHAR
			ORDER BY p.uuid FOR UPDATE;`
		args = append(args, string(filter.PartType))
	default:
		query += `ORDER BY p.uuid FOR UPDATE;`
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query, args...)
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

	return repoConverter.PartsRecordToModel(parts)
}
