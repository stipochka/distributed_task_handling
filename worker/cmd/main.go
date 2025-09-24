package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"worker/internal/consumer"
	"worker/internal/pool"
	"worker/internal/storage"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("started worker")

	godotenv.Load()

	dbUrl := os.Getenv("STORAGE_PATH")
	if dbUrl == "" {
		logrus.Fatal("failed to get db_url")
	}

	brokerAddr := os.Getenv("BROKER_ADDR")
	if brokerAddr == "" {
		logrus.Fatal("failed to get broker addr")
	}

	repo, err := storage.NewPostgresStorage(dbUrl)
	if err != nil {
		logrus.WithError(err).Fatal("failed to init repo")
	}

	pool := pool.NewPool(5, repo)

	broker, err := consumer.NewConsumer(brokerAddr, pool)
	if err != nil {
		logrus.WithError(err).Fatal("failed to initialize broker")
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	broker.Run(ctx, "tasks")

	<-done

	logrus.Info("stopping worker...")
	cancel()
}
