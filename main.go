package main

import (
	"context"
	"flag"
	"log"
	"mailer/config"
	"mailer/internal/router"
	"mailer/internal/sender"
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
		cfg           = config.ReadConfigFromFile(confPath)
		db            = mongo.New(ctx, cfg.Mongo)
		loggerConn    = rabbit.NewConn(ctx, cfg.Rabbit.Logs.Url).Publisher(cfg.Rabbit.Logs.QueueName)
		emailConsumer = rabbit.NewConn(ctx, cfg.Rabbit.Mails.Url).Consumer(ctx, cfg.Rabbit.Mails.QueueName)
		sending       = sender.New(cfg.Email)
	)

	// --------------- can't fail ---------------
	var (
		clogger = clog.New(loggerConn, cfg.Server.Name)
		routing = router.New(
			clogger,
			router.NewRepo(db.Collection("templates")),
			sending,
			emailConsumer,
		)
	)

	routing.ProcessEmails()
}
