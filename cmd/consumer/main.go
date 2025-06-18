package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"kafka-go/internal/app/consumer"
	"kafka-go/internal/config"
	"kafka-go/internal/infrastructure/logger"
)

func main() {
	log := logger.Load()
	defer log.Sync()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := logger.WithLogger(context.Background(), log)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := consumer.Run(ctx, *cfg)
		if err != nil {
			log.Errorw("failed to run consumer", "error", err)
			sigChan <- syscall.SIGTERM
		}
	}()

	log.Infow("received signal, gracefully shutting down", "signal", <-sigChan)
	cancel()
	log.Info("shutdown complete")
}
