package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/IBM/sarama"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/config"
	orderpaidconsumer "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/consumer/order_paid"
	shipassembledproducer "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/producer/ship_assembled"
	assembledsvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/service/assembly"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/closer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/consumer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/producer"
	kafkamiddleware "github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/middleware/kafka"
)

type OrderPaidConsumer interface {
	RunConsumer(ctx context.Context) error
}

type diContainer struct {
	// Инфраструктура.
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup

	// Обёртки.
	wrappedShipAssembledProducer *producer.Producer
	wrappedOrderPaidConsumer     *consumer.Consumer

	// Сервисы.
	shipAssembledProducer assembledsvc.ShipAssembledProducerService
	orderPaidConsumer     OrderPaidConsumer
	assembleSvc           orderpaidconsumer.AssembleService

	// Вспомогательные утилиты.
	sleeper assembledsvc.Sleeper
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().ShipAssembledProducer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать sync producer", "error", err.Error())
			os.Exit(1)
		}

		closer.Add("kafka sync producer", func(_ context.Context) error { return p.Close() })

		slog.Info("🚀 Kafka-продюсер запущен")

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) WrappedShipAssembledProducer() *producer.Producer {
	if d.wrappedShipAssembledProducer == nil {
		d.wrappedShipAssembledProducer = producer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().ShipAssembledProducer.Topic(),
		)
	}

	return d.wrappedShipAssembledProducer
}

func (d *diContainer) ShipAssembledProducerService() assembledsvc.ShipAssembledProducerService {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = shipassembledproducer.New(d.WrappedShipAssembledProducer())
	}

	return d.shipAssembledProducer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		cg, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers,
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.SaramaConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать consumer group", "error", err)
			os.Exit(1)
		}

		d.consumerGroup = cg
	}

	return d.consumerGroup
}

func (d *diContainer) WrappedOrderPaidConsumer() *consumer.Consumer {
	if d.wrappedOrderPaidConsumer == nil {
		d.wrappedOrderPaidConsumer = consumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			consumer.WithMiddlewares(
				kafkamiddleware.ConsumerLogging(),
			),
		)
	}

	return d.wrappedOrderPaidConsumer
}

func (d *diContainer) OrderPaidConsumer() OrderPaidConsumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = orderpaidconsumer.New(d.WrappedOrderPaidConsumer(), d.AssembleService())
	}

	return d.orderPaidConsumer
}

func (d *diContainer) AssembleService() orderpaidconsumer.AssembleService {
	if d.assembleSvc == nil {
		d.assembleSvc = assembledsvc.New(d.ShipAssembledProducerService(), d.Sleeper())
	}

	return d.assembleSvc
}

func (d *diContainer) Sleeper() assembledsvc.Sleeper {
	if d.sleeper == nil {
		d.sleeper = assembledsvc.NewSleeper(config.AppConfig().Assembler.AssembleLimitTimeSec)
	}

	return d.sleeper
}
