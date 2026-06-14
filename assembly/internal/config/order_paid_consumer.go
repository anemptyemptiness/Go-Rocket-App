package config

import "github.com/IBM/sarama"

type orderPaidConsumerConfig struct {
	TopicName string `yaml:"topic" env:"ORDER_PAID_TOPIC_NAME" env-default:"order.paid"`
	Group     string `yaml:"group_id" env:"ORDER_PAID_GROUP_ID" env-default:"1"`
}

func (o *orderPaidConsumerConfig) Topic() string { return o.TopicName }

func (o *orderPaidConsumerConfig) GroupID() string { return o.Group }

func (o *orderPaidConsumerConfig) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return cfg
}
