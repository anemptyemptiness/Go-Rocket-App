package order

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	const queryOrder = `
		SELECT
			o.uuid,
			o.total_price,
			o."status",
			o.transaction_uuid,
			o.payment_method,
			o.created_at,
			o.updated_at
		FROM public.orders AS o
		WHERE o.uuid = $1;`

	var order record.Order

	err := r.getter.DefaultTrOrDB(ctx, r.pool).QueryRow(ctx, queryOrder, orderUUID).Scan(
		&order.Uuid,
		&order.TotalPrice,
		&order.Status,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return model.Order{}, handleError(err)
	}

	const queryItems = `
		SELECT
			uuid,
			order_uuid,
			part_uuid,
			part_type,
			price,
			created_at
		FROM order_items
		WHERE order_uuid = $1
		ORDER BY created_at;`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, queryItems, orderUUID)
	if err != nil {
		return model.Order{}, err
	}
	defer rows.Close()

	var items []record.OrderItem
	for rows.Next() {
		var item record.OrderItem
		if err := rows.Scan(
			&item.Uuid,
			&item.OrderUuid,
			&item.PartUuid,
			&item.PartType,
			&item.Price,
			&item.CreatedAt,
		); err != nil {
			return model.Order{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return model.Order{}, err
	}

	order.OrderItems = items

	return repoConverter.OrderRecordToModel(order), nil
}

func handleError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errs.ErrOrderNotFound
	default:
		return err
	}
}
