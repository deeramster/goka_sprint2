package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/deeramster/goka_sprint2/pkg/censor"
	"github.com/deeramster/goka_sprint2/pkg/kafka"
	"github.com/deeramster/goka_sprint2/pkg/processor"
	"github.com/deeramster/goka_sprint2/pkg/storage"
)

func main() {
	storageImpl, err := storage.NewFileStorage("./data/blocked_users")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	censorService := censor.NewCensor([]string{"bad", "words", "censored"})

	msgProcessor := processor.NewMessageProcessor(storageImpl, censorService)

	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		brokersStr = "localhost:9094,localhost:9095,localhost:9096"
	}
	brokers := strings.Split(brokersStr, ",")

	kafkaProcessor := kafka.NewKafkaProcessor(brokers, msgProcessor)

	if err := kafkaProcessor.Run(context.Background()); err != nil {
		log.Fatalf("Error running processor: %v", err)
	}
}
