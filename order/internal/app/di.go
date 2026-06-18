package app

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/IBM/sarama"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderapi "github.com/anemptyemptiness/Go-Rocket-App/order/internal/api/order/v1"
	inventoryclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/inventory/v1"
	paymentclientv1 "github.com/anemptyemptiness/Go-Rocket-App/order/internal/client/grpc/payment/v1"
	"github.com/anemptyemptiness/Go-Rocket-App/order/internal/config"
	shipassembledconsumer "github.com/anemptyemptiness/Go-Rocket-App/order/internal/consumer/assembly_consumer"
	orderpaidproducer "github.com/anemptyemptiness/Go-Rocket-App/order/internal/producer/order_producer"
	orderrepo "github.com/anemptyemptiness/Go-Rocket-App/order/internal/repository/order"
	ordersvc "github.com/anemptyemptiness/Go-Rocket-App/order/internal/service/order"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/consumer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/producer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/middleware/kafka"
	orderv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/anemptyemptiness/Go-Rocket-App/shared/pkg/proto/payment/v1"
)

type ShipAssembledConsumer interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	orderHandler    http.Handler
	orderAPI        orderv1.Handler
	orderService    orderapi.OrderService
	orderRepository ordersvc.OrderRepository
	txManager       ordersvc.TxManager
	paymentClient   ordersvc.PaymentClient
	inventoryClient ordersvc.InventoryClient
	pool            *pgxpool.Pool

	// Кафка: обёртки, продюсеры, консюмеры, сервисы.
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup

	wrapperShipAssembledConsumer *consumer.Consumer
	wrapperOrderPaidProducer     *producer.Producer

	shipAssembledConsumer ShipAssembledConsumer
	orderPaidProducer     ordersvc.OrderPaidProducerService
}

func (d *diContainer) OrderServer(ctx context.Context) http.Handler {
	if d.orderHandler == nil {
		server, err := orderv1.NewServer(d.OrderAPI(ctx))
		if err != nil {
			slog.Error("создание http хэндлера", "error", err)
			os.Exit(1)
		}

		d.orderHandler = server

		slog.Info("http хэндлер успешно создан")
	}

	return d.orderHandler
}

func (d *diContainer) OrderAPI(ctx context.Context) orderv1.Handler {
	if d.orderAPI == nil {
		d.orderAPI = orderapi.NewAPI(d.OrderService(ctx))
	}

	return d.orderAPI
}

func (d *diContainer) OrderService(ctx context.Context) orderapi.OrderService {
	if d.orderService == nil {
		d.orderService = ordersvc.New(
			d.OrderRepository(ctx),
			d.PaymentClient(ctx),
			d.InventoryClient(ctx),
			d.TxManager(ctx),
			d.OrderPaidProducer(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) ordersvc.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderrepo.New(d.PGPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) TxManager(ctx context.Context) ordersvc.TxManager {
	if d.txManager == nil {
		txManager, err := manager.New(trmpgx.NewDefaultFactory(d.PGPool(ctx)))
		if err != nil {
			slog.Error("создание transaction manager", "error", err)
			os.Exit(1)
		}

		d.txManager = txManager

		slog.Info("transaction manager успешно создан")
	}

	return d.txManager
}

//nolint:dupl // Два разных gRPC-клиента с одинаковым шаблоном инициализации; намеренно оставлено явно.
func (d *diContainer) PaymentClient(_ context.Context) ordersvc.PaymentClient {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(config.AppConfig().GRPC.PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                keepAliveTime,
				Timeout:             keepAliveTimeout,
				PermitWithoutStream: keepAlivePermitWithoutStream,
			}),
		)
		if err != nil {
			slog.Error("не удалось подключиться к PaymentService", "error", err)
			os.Exit(1)
		}

		closer.Add("payment GRPC Client", func(_ context.Context) error {
			err = paymentConn.Close()
			if err != nil {
				slog.Error("не удалось закрыть gRPC соединение payment client", "error", err)
			} else {
				slog.Info("соединение с PaymentClient закрыто")
			}
			return nil
		})

		paymentClientGRPC := paymentv1.NewPaymentServiceClient(paymentConn)
		d.paymentClient = paymentclientv1.New(paymentClientGRPC)

		slog.Info("подключение к PaymentClient установлено")
	}

	return d.paymentClient
}

//nolint:dupl // Два разных gRPC-клиента с одинаковым шаблоном инициализации; намеренно оставлено явно.
func (d *diContainer) InventoryClient(_ context.Context) ordersvc.InventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(config.AppConfig().GRPC.InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                keepAliveTime,
				Timeout:             keepAliveTimeout,
				PermitWithoutStream: keepAlivePermitWithoutStream,
			}),
		)
		if err != nil {
			slog.Error("не удалось подключиться к InventoryService", "error", err)
			os.Exit(1)
		}

		closer.Add("inventory GRPC Client", func(_ context.Context) error {
			err = inventoryConn.Close()
			if err != nil {
				slog.Error("не удалось закрыть gRPC соединение inventory client", "error", err)
			} else {
				slog.Info("соединение с InventoryClient закрыто")
			}
			return nil
		})

		inventoryClientGRPC := inventoryv1.NewInventoryServiceClient(inventoryConn)
		d.inventoryClient = inventoryclientv1.New(inventoryClientGRPC)

		slog.Info("подключение к InventoryClient установлено")
	}

	return d.inventoryClient
}

func (d *diContainer) PGPool(ctx context.Context) *pgxpool.Pool {
	if d.pool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().PG.DSN())
		if err != nil {
			slog.Error("создание пула соединений", "error", err)
			os.Exit(1)
		}

		closer.Add("pgPool", func(_ context.Context) error {
			pool.Close()
			slog.Info("соединение с PostgreSQL закрыто")
			return nil
		})

		err = pool.Ping(ctx)
		if err != nil {
			slog.Error("проверка соединения с БД", "error", err)
			os.Exit(1)
		}

		slog.Info("подключение к PostgreSQL установлено")

		d.pool = pool
	}

	return d.pool
}

func (d *diContainer) ConsumerGroup(_ context.Context) sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cg, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().ShipAssembledConsumer.GroupID(),
			config.AppConfig().ShipAssembledConsumer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать kafka consumer group", "error", err)
			os.Exit(1)
		}

		d.consumerGroup = cg
	}

	return d.consumerGroup
}

func (d *diContainer) WrappedShipAssembledConsumer(ctx context.Context) *consumer.Consumer {
	if d.wrapperShipAssembledConsumer == nil {
		d.wrapperShipAssembledConsumer = consumer.NewConsumer(
			d.ConsumerGroup(ctx),
			[]string{
				config.AppConfig().ShipAssembledConsumer.Topic(),
			},
			consumer.WithMiddlewares(kafka.ConsumerLogging()),
		)
	}

	return d.wrapperShipAssembledConsumer
}

func (d *diContainer) ShipAssembledConsumer(ctx context.Context) ShipAssembledConsumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = shipassembledconsumer.New(
			d.WrappedShipAssembledConsumer(ctx),
			d.OrderService(ctx),
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
		)
	}

	return d.shipAssembledConsumer
}

func (d *diContainer) SyncProducer(_ context.Context) sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().OrderPaidProducer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать kafka sync producer", "error", err)
			os.Exit(1)
		}

		closer.Add("kafka sync producer", func(_ context.Context) error { return p.Close() })

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) WrappedOrderPaidProducer(ctx context.Context) *producer.Producer {
	if d.wrapperOrderPaidProducer == nil {
		d.wrapperOrderPaidProducer = producer.NewProducer(
			d.SyncProducer(ctx),
			config.AppConfig().OrderPaidProducer.Topic(),
		)
	}

	return d.wrapperOrderPaidProducer
}

func (d *diContainer) OrderPaidProducer(ctx context.Context) ordersvc.OrderPaidProducerService {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = orderpaidproducer.New(d.WrappedOrderPaidProducer(ctx))
	}

	return d.orderPaidProducer
}
