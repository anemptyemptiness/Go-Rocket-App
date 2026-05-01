package part

import (
	"context"

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
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return model.Part{}, errs.ErrPartNotFound
		default:
			return model.Part{}, err
		}
	}

	return repoConverter.PartRecordToModel(part), nil
}

func (r *repository) ListParts(ctx context.Context, filter model.PartFilter) ([]model.Part, error) {
	var rows pgx.Rows
	var err error

	const query = `
		SELECT
			p.uuid,
			p.name,
			p.description,
			p.part_type,
			p.price,
			p.stock_quantity,
			p.created_at
		FROM parts AS p;`

	rows, err = r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, query)
	if err != nil {
		return nil, handleError(err)
	}
	defer rows.Close()

	modelParts, err := scanParts(rows)
	if err != nil {
		return nil, err
	}

	return modelParts, nil
}

func scanParts(rows pgx.Rows) ([]model.Part, error) {
	var parts []record.Part

	for rows.Next() {
		part := record.Part{}

		if scanErr := rows.Scan(
			&part.UUID,
			&part.Name,
			&part.Description,
			&part.PartType,
			&part.Price,
			&part.StockQuantity,
			&part.CreatedAt,
		); scanErr != nil {
			return nil, scanErr
		}

		parts = append(parts, part)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return repoConverter.PartsRecordToModel(parts), nil
}

func handleError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.ErrPartNotFound
	default:
		return err
	}
}
