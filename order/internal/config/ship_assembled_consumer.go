package config

import "github.com/IBM/sarama"

type shipAssembledConsumer struct {
	TopicName string `yaml:"topic" env-default:"assembly.ship-assembled"`
	Group     string `yaml:"group_id" env-default:"2"`
}

func (o *shipAssembledConsumer) Topic() string { return o.TopicName }

func (o *shipAssembledConsumer) GroupID() string { return o.Group }

func (o *shipAssembledConsumer) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return cfg
}
