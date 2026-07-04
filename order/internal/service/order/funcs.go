package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	errs "github.com/anemptyemptiness/Go-Rocket-App/order/internal/errors"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/model"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/input"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/auth"
	pkgerr "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/errors"
)

func (s *service) Get(ctx context.Context, orderUUID string) (model.Order, error) {
	userUUID, ok := auth.UserUUIDFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "получение заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrEmptyUserUUID.Error()),
		)

		return model.Order{}, pkgerr.Unauthenticated(errs.ErrEmptyUserUUID)
	}

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		slog.ErrorContext(ctx, "получение заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", err.Error()),
		)

		if errors.Is(err, errs.ErrOrderNotFound) {
			return model.Order{}, pkgerr.NotFound(err)
		}
		return model.Order{}, pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
	}

	return order, nil
}

func (s *service) Create(ctx context.Context, req input.CreateOrderRequest) (model.Order, error) {
	userUUID, ok := auth.UserUUIDFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "создание заказа",
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrEmptyUserUUID.Error()),
		)
		return model.Order{}, pkgerr.Unauthenticated(errs.ErrEmptyUserUUID)
	}

	seen := make(map[string]struct{})
	for _, id := range req.PartUUIDs() {
		if id == "" {
			slog.ErrorContext(ctx, "создание заказа",
				slog.String("order_uuid", id),
				slog.String("error", errs.ErrInvalidPartUUID.Error()),
			)

			return model.Order{}, pkgerr.InvalidArgument(errs.ErrInvalidPartUUID)
		}
		if _, exists := seen[id]; exists {
			slog.ErrorContext(ctx, "создание заказа",
				slog.String("order_uuid", id),
				slog.String("error", errs.ErrInvalidPartUUID.Error()),
			)

			return model.Order{}, pkgerr.InvalidArgument(errs.ErrInvalidPartUUID)
		}

		seen[id] = struct{}{}
	}

	var order model.Order

	err := s.txManager.Do(ctx, func(txCtx context.Context) error {
		listPartsCtx, cancel := context.WithTimeout(txCtx, 5*time.Second)
		defer cancel()

		parts, err := s.inventoryClient.ListParts(listPartsCtx, req.PartUUIDs())
		if err != nil {
			return err
		}

		var totalPrice int64
		var orderItems []model.OrderItem
		for _, part := range parts {
			if part.StockQuantity <= 0 {
				return pkgerr.Conflict(fmt.Errorf("деталь %s: %w", part.Name, errs.ErrPartIsOver))
			}

			orderItems = append(orderItems, model.OrderItem{
				PartUuid: part.UUID,
				PartType: part.PartType,
				Price:    part.Price,
			})

			totalPrice += part.Price
		}

		order.UserUUID = userUUID.String()
		order.Items = orderItems
		order.TotalPrice = totalPrice
		order.Status = model.OrderStatusPendingPayment

		validatePartsCtx, cancel := context.WithTimeout(txCtx, 5*time.Second)
		defer cancel()

		err = s.inventoryClient.ValidateCompatibility(validatePartsCtx, req)
		if err != nil {
			return err
		}

		reservePartsCtx, cancel := context.WithTimeout(txCtx, 5*time.Second)
		defer cancel()

		err = s.inventoryClient.ReserveParts(reservePartsCtx, req.PartUUIDs())
		if err != nil {
			return err
		}

		orderUUID, err := s.orderRepository.Create(txCtx, order)
		if err != nil {
			return pkgerr.Internal(fmt.Errorf("создать заказ: %w", err))
		}

		order.SetID(orderUUID)

		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "создание заказа", slog.String("error", err.Error()))
		return model.Order{}, err
	}

	ordersCreatedTotal.Add(ctx, 1)

	return order, nil
}

func (s *service) Pay(ctx context.Context, orderUUID string, method model.PaymentMethod) (string, error) {
	var transactionUUID string

	userUUID, ok := auth.UserUUIDFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "оплата заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrEmptyUserUUID.Error()),
		)
		return "", pkgerr.Unauthenticated(errs.ErrEmptyUserUUID)
	}

	var order model.Order
	var err error

	err = s.txManager.Do(ctx, func(txCtx context.Context) error {
		order, err = s.orderRepository.GetForUpdate(ctx, orderUUID)
		if err != nil {
			if errors.Is(err, errs.ErrOrderNotFound) {
				return pkgerr.NotFound(err)
			}
			return pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
		}

		switch order.Status {
		case model.OrderStatusPaid:
			return pkgerr.Conflict(errs.ErrOrderAlreadyPaid)
		case model.OrderStatusCancelled:
			return pkgerr.Conflict(errs.ErrOrderAlreadyCancelled)
		case model.OrderStatusAssembled:
			return pkgerr.Conflict(errs.ErrOrderAssembled)
		}

		if order.Status != model.OrderStatusPendingPayment {
			return pkgerr.InvalidArgument(errs.ErrOrderStatusIncorrect)
		}

		clientCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		transactionUUID, err = s.paymentClient.PayOrder(clientCtx, orderUUID, method)
		if err != nil {
			return fmt.Errorf("оплатить заказ: %w", err)
		}

		order.SetStatus(model.OrderStatusPaid)
		order.SetTransactionID(transactionUUID)
		order.SetPaymentMethod(method)

		err = s.orderRepository.Update(ctx, order)
		if err != nil {
			if errors.Is(err, errs.ErrOrderNotFound) {
				return pkgerr.NotFound(err)
			}
			return pkgerr.Internal(fmt.Errorf("обновить заказ: %w", err))
		}

		eventUUID := uuid.New().String()

		err = s.orderPaidProducer.Produce(ctx, model.NewOrderPaidEvent(
			eventUUID,
			order.UUID,
			order.UserUUID,
		))
		if err != nil {
			return pkgerr.Internal(fmt.Errorf("отправка ивента %s OrderPaid в брокер сообщений: %w", eventUUID, err))
		}

		slog.InfoContext(ctx, "ивент OrderPaid успешно отправлен", "event_uuid", eventUUID)

		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, "оплата заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("payment_method", string(method)),
			slog.String("error", err.Error()),
		)
		return "", err
	}

	ordersPaidTotal.Add(ctx, 1)
	ordersRevenueTotal.Add(ctx, order.TotalPrice)

	return transactionUUID, nil
}

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	userUUID, ok := auth.UserUUIDFromContext(ctx)
	if !ok {
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrEmptyUserUUID.Error()),
		)
		return pkgerr.Unauthenticated(errs.ErrEmptyUserUUID)
	}

	order, err := s.orderRepository.GetForUpdate(ctx, orderUUID)
	if err != nil {
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", err.Error()),
		)

		if errors.Is(err, errs.ErrOrderNotFound) {
			return pkgerr.NotFound(err)
		}
		return pkgerr.Internal(fmt.Errorf("получить заказ: %w", err))
	}

	switch order.Status {
	case model.OrderStatusPaid:
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrOrderAlreadyPaid.Error()),
		)
		return pkgerr.Conflict(errs.ErrOrderAlreadyPaid)
	case model.OrderStatusCancelled:
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrOrderAlreadyCancelled.Error()),
		)
		return pkgerr.Conflict(errs.ErrOrderAlreadyCancelled)
	case model.OrderStatusAssembled:
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", errs.ErrOrderAssembled.Error()),
		)
		return pkgerr.Conflict(errs.ErrOrderAssembled)
	}

	uuids := make([]string, 0, len(order.Items))
	for _, item := range order.Items {
		uuids = append(uuids, item.PartUuid)
	}

	releasePartsCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// На будущее: здесь возможно реализовать паттерн SAGA, чтобы отменить ReleaseParts в grpc-inventory.
	err = s.inventoryClient.ReleaseParts(releasePartsCtx, uuids)
	if err != nil {
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", err.Error()),
		)
		return err
	}

	order.SetStatus(model.OrderStatusCancelled)

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		slog.ErrorContext(ctx, "отмена заказа",
			slog.String("order_uuid", orderUUID),
			slog.String("user_uuid", userUUID.String()),
			slog.String("error", err.Error()),
		)

		if errors.Is(err, errs.ErrOrderNotFound) {
			return pkgerr.NotFound(err)
		}
		return pkgerr.Internal(fmt.Errorf("обновить заказ: %w", err))
	}

	return nil
}
