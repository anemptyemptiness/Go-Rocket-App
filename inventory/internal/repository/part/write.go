package part

import (
	"context"

	"github.com/anemptyemptiness/Go-Rocket-App/inventory/internal/model/entity"
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
