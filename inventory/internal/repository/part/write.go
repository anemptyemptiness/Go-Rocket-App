package part

import (
	"context"

	errs "github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/service/input"
)

func (r *repository) UpdateReservedBatch(ctx context.Context, parts []entity.Part) error {
	const query = `
		UPDATE parts AS p
		SET
			reserved = batch.reserved,
			updated_at = NOW()
		FROM unnest($1::UUID[], $2::INT[]) AS batch(uuid, reserved)
		WHERE p.uuid = batch.uuid;`

	uuids := make([]string, len(parts))
	reserved := make([]int64, len(parts))

	for idx, part := range parts {
		uuids[idx] = part.GetPartUUID()
		reserved[idx] = part.GetReserved()
	}

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, uuids, reserved)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CommitParts(ctx context.Context, req input.CommitPartsRequest) error {
	const query = `
		UPDATE parts
		SET
		    stock_quantity = stock_quantity - 1,
		    reserved = reserved - 1
		WHERE uuid = ANY($1::UUID[])
			AND stock_quantity > 0
			AND reserved > 0;`

	cmdTag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, req.UUIDs)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != int64(len(req.UUIDs)) {
		return errs.ErrPartNotFound
	}

	return nil
}
