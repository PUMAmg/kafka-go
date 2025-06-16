package main

import (
	"kafka-go/internal/infrastructure/logger"

	"github.com/joho/godotenv"
)

func main() {
	log := logger.Load()
	defer log.Sync()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env file %v", err)
	}

	//cfg, err := config.Load()

}
