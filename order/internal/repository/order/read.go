package order

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	repoConverter "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/converter"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := r.get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	orderItems, err := r.getOrderItems(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	order.OrderItems = orderItems

	return repoConverter.OrderRecordToModel(order), nil
}

func (r *repository) get(ctx context.Context, orderUUID string) (record.Order, error) {
	const queryOrder = `
		SELECT
			o.uuid,
			o.user_uuid,
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
		&order.UserUuid,
		&order.TotalPrice,
		&order.Status,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return record.Order{}, fmt.Errorf("%w: uuid=%s", errs.ErrOrderNotFound, orderUUID)
		}
		return record.Order{}, err
	}

	return order, nil
}

func (r *repository) getOrderItems(ctx context.Context, orderUUID string) ([]record.OrderItem, error) {
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
		return nil, err
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[record.OrderItem])
	if err != nil {
		return nil, err
	}

	return items, nil
}
