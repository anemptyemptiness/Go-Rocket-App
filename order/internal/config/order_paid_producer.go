package config

import "github.com/IBM/sarama"

type orderPaidProducer struct {
	TopicName string `yaml:"topic" env-default:"order.paid"`
}

func (c *orderPaidProducer) Topic() string {
	return c.TopicName
}

func (c *orderPaidProducer) SaramaConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V4_0_0_0
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5

	return cfg
}
