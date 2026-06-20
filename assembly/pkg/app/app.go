package app

import (
	"github.com/IBM/sarama"

	"github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/app"
	consumersvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/consumer/order_paid"
	producersvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/producer/ship_assembled"
	assemblysvc "github.com/anemptyemptiness/Go-Rocket-App/assembly/internal/service/assembly"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/consumer"
	"github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/kafka/producer"
	kafkamiddleware "github.com/anemptyemptiness/Go-Rocket-App/platform/pkg/middleware/kafka"
)

type Config struct {
	OrderPaidTopic     string
	ShipAssembledTopic string
}

func New(
	syncProducer sarama.SyncProducer,
	consumerGroup sarama.ConsumerGroup,
	cfg Config,
) app.OrderPaidConsumer {
	wrapperConsumer := consumer.NewConsumer(
		consumerGroup,
		[]string{cfg.OrderPaidTopic},
		consumer.WithMiddlewares(
			kafkamiddleware.ConsumerLogging(),
			kafkamiddleware.ConsumerSession(),
		),
	)

	wrapperProducer := producer.NewProducer(
		syncProducer,
		cfg.ShipAssembledTopic,
	)

	producerSvc := producersvc.New(wrapperProducer)
	sleeper := assemblysvc.NewSleeper(1)
	assemblySvc := assemblysvc.New(producerSvc, sleeper)

	svc := consumersvc.New(wrapperConsumer, assemblySvc)

	return svc
}
