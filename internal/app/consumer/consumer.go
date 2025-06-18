package consumer

import (
	"context"

	"kafka-go/internal/config"
	"kafka-go/internal/infrastructure/logger"
	"kafka-go/internal/kafka"
)

const appName = "kafka-consumer"

func Run(ctx context.Context, cfg config.Config) error {
	const op = "consumer.Run"

	log := logger.FromContext(ctx)

	consumer, err := kafka.NewConsumer(cfg, appName)
	if err != nil {
		log.Errorf("%s: failed to create consumer %v", op, err)
		return err
	}
	defer consumer.Close()

	err = consumer.Consume(ctx, func(msg []byte) error {
		log.Infof("processing message: %s", string(msg))
		return nil
	})
	if err != nil {
		log.Errorf("%s: failed to consume messages %v", op, err)
		return err
	}

	<-ctx.Done()
	log.Info("shutting down")

	return nil
}
