package service

import (
	"context"
	"encoding/json"
	"fmt"

	"kafka-go/internal/domain"
	"kafka-go/internal/infrastructure/logger"

	"github.com/brianvoe/gofakeit/v7"
)

type Producer interface {
	Produce(ctx context.Context, message []byte) error
}

type MessageService struct {
	messageProducer Producer
}

func NewMessageService(messageProducer Producer) *MessageService {
	return &MessageService{
		messageProducer: messageProducer,
	}
}

func (s *MessageService) SendMessage(ctx context.Context) error {
	const op = "service.SendMessage"

	log := logger.FromContext(ctx)

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := generateMessage()
			if err != nil {
				log.Errorf("failed to generate message %v", err)
				return fmt.Errorf("%s:%s", op, err)
			}

			if err = s.messageProducer.Produce(ctx, msg); err != nil {
				log.Errorf("failed to produce message %v", err)
				return fmt.Errorf("%s:%s", op, err)
			}
		}
	}
	return nil
}

func generateMessage() ([]byte, error) {
	const op = "service.generateMessage"
	text := gofakeit.ProductDescription()

	message := domain.Message{Text: text}

	rawMsg, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to marshal message %v", op, err)
	}
	return rawMsg, nil
}
