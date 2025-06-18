package producer

import (
	"context"

	"kafka-go/internal/config"
	"kafka-go/internal/infrastructure/logger"
	"kafka-go/internal/kafka"
	"kafka-go/internal/service"
)

const appName = "kafka-producer"

func Run(ctx context.Context, cfg config.Config) error {
	const op = "producer.Run"

	log := logger.FromContext(ctx)

	producer, err := kafka.NewProducer(cfg, appName)
	if err != nil {
		log.Errorf("%s: failed to create producer %v", op, err)
		return err
	}
	defer producer.Close()

	srv := service.NewMessageService(producer)

	if err := srv.SendMessage(ctx); err != nil {
		log.Errorf("%s: failed to send message %v", op, err)
	}

	<-ctx.Done()
	log.Info("shutting down")

	return nil
}
