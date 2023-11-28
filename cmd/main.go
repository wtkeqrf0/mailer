package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"log"
	"mailer/config"
	consumerRepository "mailer/internal/consumer/repository"
	"mailer/internal/consumer/usecase"
	singleSenderUseCase "mailer/internal/single_sender/usecase"
	"mailer/pkg/guzzle_logger"
	"mailer/pkg/logger"
	"mailer/pkg/mongodb"
	rabbitConsumer "mailer/pkg/rabbit/consumer"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config.yml: %v", err)
	}

	apiLogger, err := logger.NewApiLogger(cfg.Logger, cfg.ServiceName)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	mongoClient, err := mongoDB.New(cfg.Mongo)
	if err != nil {
		apiLogger.Fatalf("failed to init mongo connection: %v", err)
	}
	const collectionTemplates = "templates"

	guzLog, err := guzzle_logger.New(cfg.ServiceName, "email_single_sender", apiLogger, cfg.Rabbit.GuzzleLogger)
	if err != nil {
		apiLogger.Fatalf("failed to init guzzle publisher: %v", err)
	}

	// singleSenderRepo := singleSenderRepository.NewSingleSenderRepo(mongoClient.Collection(collectionTemplates))
	singleSenderUC, err := singleSenderUseCase.NewSingleSender(cfg.SingleSender)
	if err != nil {
		apiLogger.Fatalf("failed to init email single sender connection: %v", err)
	}

	consumerRepo := consumerRepository.NewConsumerRepo(mongoClient.Collection(collectionTemplates))
	consumerUC := usecase.NewConsumer(apiLogger, guzLog, consumerRepo, singleSenderUC)

	consumer, err := rabbitConsumer.StartConsuming(
		cfg.Rabbit.Consumer,
		consumerUC.ProcessEmail,
	)
	if err != nil {
		apiLogger.Fatalf("failed to start consuming: %v", err)
	}
	defer consumer.Close()
	apiLogger.Info("awaiting signal...")

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	apiLogger.Infof("received signal %v. Shutdown...", <-sigs)
}
