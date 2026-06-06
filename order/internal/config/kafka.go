package config

type kafkaConfig struct {
	ExternalPort   string      `yaml:"external_port" env:"KAFKA_EXTERNAL_PORT" env-default:"9092"`
	InternalPort   string      `yaml:"internal_port" env:"KAFKA_INTERNAL_PORT" env-default:"29092"`
	ControllerPort string      `yaml:"controller_port" env:"KAFKA_CONTROLLER_PORT" env-default:"29093"`
	UIPort         string      `yaml:"ui_port" env:"KAFKA_UI_PORT" env-default:"8090"`
	Topic          topicConfig `yaml:"topic"`
}

type topicConfig struct {
	Pay      string `yaml:"pay" env:"KAFKA_TOPIC_PAY" env-default:"order.paid"`
	Assembly string `yaml:"assembly" env:"KAFKA_TOPIC_ASSEMBLY" env-default:"assembly.ship-assembled"`
}
