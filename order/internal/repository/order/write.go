package order

import (
	"context"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) (string, error) {
	orderUUID, err := r.createOrder(ctx, order)
	if err != nil {
		return "", err
	}

	order.UUID = orderUUID

	err = r.createOrderItems(ctx, order)
	if err != nil {
		return "", err
	}

	return order.UUID, nil
}

func (r *repository) createOrder(ctx context.Context, order model.Order) (string, error) {
	orderUUID := uuid.New()

	builder := sqlbuilder.NewInsertBuilder()
	builder.InsertInto("orders").Cols("uuid", "user_uuid", "total_price", "status", "payment_method")
	builder.Values(orderUUID, order.UserUUID, order.TotalPrice, order.Status, order.PaymentMethod)

	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)

	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx, query,
		args...,
	)
	if err != nil {
		return "", err
	}

	return orderUUID.String(), nil
}

func (r *repository) createOrderItems(ctx context.Context, order model.Order) error {
	builder := sqlbuilder.NewInsertBuilder()
	builder.InsertInto("order_items").Cols("uuid", "order_uuid", "part_uuid", "part_type", "price")

	for _, item := range order.Items {
		builder.Values(
			uuid.New(),
			order.UUID,
			item.PartUuid,
			item.PartType,
			item.Price,
		)
	}

	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	_, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, order model.Order) error {
	const query = `
		UPDATE orders
		SET
		    total_price = COALESCE($2, total_price),
		    "status" = COALESCE($3, "status"),
		    transaction_uuid = COALESCE($4, transaction_uuid),
		    payment_method = COALESCE($5, payment_method),
		    updated_at = NOW()
		WHERE uuid = $1;`

	tag, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx, query,
		order.UUID,
		order.TotalPrice,
		order.Status,
		order.TransactionUUID,
		order.PaymentMethod,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrOrderNotFound
	}

	return nil
}
