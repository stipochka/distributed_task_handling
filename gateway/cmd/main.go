package main

import (
	"context"
	"gateway/internal/handler"
	"gateway/internal/producer"
	"gateway/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

const (
	addr            = "127.0.0.1:8080"
	shutdownTimeout = time.Second * 5
)

func main() {

	log.WithField(
		"time", time.Now,
	).Info("Started api gateway")

	err := godotenv.Load()
	if err != nil {
		log.WithField(
			"error", err.Error(),
		).Fatal("failed to load .env")
	}

	natsAddr := os.Getenv("NATS_ADDR")
	if natsAddr == "" {
		log.Fatal("empty nats addr")
	}

	prod, err := producer.NewNatsProducer(natsAddr)
	if err != nil {
		log.WithField(
			"error", err.Error(),
		).Fatal("failed to connect to broker")
	}

	useCase := service.NewTaskSender(prod)

	handler := handler.NewHandler(useCase)

	router := handler.InitRoute()

	serverAddr := os.Getenv("HTTP_SERVER")
	if serverAddr == "" {
		serverAddr = addr
	}

	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.WithField("error", err.Error()).Error("error while running server")
		}
	}()

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithField("error", err.Error()).Error("error while stopping server")

		return
	}

	log.Info("service stopped")

}
