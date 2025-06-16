package kafka

import (
	"context"
	"fmt"

	"kafka-go/internal/config"
	"kafka-go/internal/infrastructure/logger"

	"github.com/IBM/sarama"
)

type MsgProducer struct {
	producer  sarama.SyncProducer
	topicName string
}

func NewProducer(cfg config.Config, appName string) (*MsgProducer, error) {
	clientCfg := sarama.NewConfig()
	//требования для SyncProducer
	clientCfg.Producer.Return.Successes = true
	clientCfg.Producer.Return.Errors = true

	clientCfg.Producer.RequiredAcks = sarama.WaitForAll

	clientCfg.ClientID = appName
	clientCfg.Version = sarama.V4_0_0_0

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, clientCfg)
	if err != nil {
		return nil, err
	}

	return &MsgProducer{
		producer:  producer,
		topicName: cfg.Kafka.Topic,
	}, nil
}

func (p *MsgProducer) Produce(ctx context.Context, message []byte) error {
	const op = "kafka.Produce"

	log := logger.FromContext(ctx)

	msg := &sarama.ProducerMessage{
		Topic: p.topicName,
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Errorf("failed sending message %v", err)
		return fmt.Errorf("%s:%s", op, err)
	}

	log.Infof("message sent to topic %s, partition %d, offset %d", p.topicName, partition, offset)
	return nil
}

func (p *MsgProducer) Close() {
	p.producer.Close()
}
