package main

import (
	"context"
	"flag"
	"log"
	"mailer/config"
	consumerRepository "mailer/internal/consumer/repository"
	"mailer/internal/consumer/usecase"
	singleSenderUseCase "mailer/internal/single_sender/usecase"
	"mailer/pkg/clog"
	"mailer/pkg/mongo"
	"mailer/pkg/rabbit"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var confPath string
	flag.StringVar(&confPath, "config-path", "./config/config.yaml", "Path to config file")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// ---------------- may fail ----------------
	var (
		cfg           = config.ReadConfigFromFile[config.Config](confPath)
		db            = mongo.New(ctx, cfg.Mongo)
		loggerConn    = rabbit.NewConn(ctx, cfg.Rabbit.Logs.URL).Publisher(cfg.Rabbit.Logs.QueueName)
		emailConsumer = rabbit.NewConn(ctx, cfg.Rabbit.Mails.URL).Consumer(ctx, cfg.Rabbit.Mails.QueueName)
	)

	// --------------- can't fail ---------------
	clogger := clog.New(loggerConn, cfg.Server.Name)

	singleSenderUC := singleSenderUseCase.NewSingleSender(cfg.SingleSender)

	consumerRepo := consumerRepository.NewConsumerRepo(db.Collection("templates"))
	consumerUC := usecase.NewConsumer(clogger, consumerRepo, singleSenderUC, emailConsumer)

	consumerUC.ProcessEmails()
}
