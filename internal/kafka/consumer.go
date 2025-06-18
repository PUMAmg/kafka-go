package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"kafka-go/internal/config"
	"kafka-go/internal/domain"
	"kafka-go/internal/infrastructure/logger"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type MsgConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
}

func NewConsumer(cfg config.Config, appName string) (*MsgConsumer, error) {
	clientCfg := sarama.NewConfig()
	clientCfg.ClientID = appName
	clientCfg.Version = sarama.V4_0_0_0

	// Настройки для Consumer Group
	clientCfg.Consumer.Return.Errors = true
	clientCfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, appName, clientCfg)
	if err != nil {
		return nil, err
	}

	return &MsgConsumer{
		consumer: consumer,
		topic:    cfg.Kafka.Topic,
	}, nil
}

type consumerGroupHandler struct {
	handler func([]byte) error
	log     *zap.SugaredLogger
	topic   string
}

func (h *consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg.Value); err != nil {
			h.log.Errorf("failed to handle message from partition %d: %v", msg.Partition, err)
			continue
		}

		var message domain.Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			h.log.Errorf("failed to unmarshal message: %v", err)
			continue
		}
		fmt.Printf("received message: %+v\n", message)
		h.log.Infof("message consumed from topic %s, partition %d, offset %d", msg.Topic, msg.Partition, msg.Offset)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *MsgConsumer) Consume(ctx context.Context, handler func([]byte) error) error {
	const op = "kafka.Consume"

	log := logger.FromContext(ctx)
	consumerHandler := &consumerGroupHandler{
		handler: handler,
		log:     log,
		topic:   c.topic,
	}

	// Запускаем горутину для чтения ошибок
	go func() {
		for err := range c.consumer.Errors() {
			log.Errorf("%s: consumer error: %v", op, err)
		}
	}()

	topics := []string{c.topic, "test-topic"}
	for {
		err := c.consumer.Consume(ctx, topics, consumerHandler)
		if err != nil {
			log.Errorf("%s: failed to consume: %v", op, err)
			continue
		}
		if ctx.Err() != nil {
			return nil
		}
	}
}

func (c *MsgConsumer) Close() {
	c.consumer.Close()
}
