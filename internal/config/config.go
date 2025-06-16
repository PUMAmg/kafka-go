package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Kafka Kafka
}

type Kafka struct {
	Brokers []string `env:"KAFKA_BROKERS,notEmpty"`
	Topic   string   `env:"KAFKA_TOPIC,notEmpty"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
