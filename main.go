package main

import (
	"context"
	"flag"
	"log"
	"mailer/config"
	"mailer/internal/router"
	"mailer/internal/sender"
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

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	// ---------------- may fail ----------------
	var (
		cfg           = config.ReadConfigFromFile(confPath)
		db            = mongo.New(ctx, cfg.Mongo)
		emailConsumer = rabbit.NewConn(ctx, cfg.Rabbit.Email.Url).Consumer(ctx, cfg.Rabbit.Email.QueueName)
		sending       = sender.New(ctx, cfg.Email)
	)

	// --------------- can't fail ---------------
	var (
		routing = router.New(
			router.NewRepo(db.Collection("templates")),
			sending,
			emailConsumer,
		)
	)

	log.Println("Service started successfully")
	routing.ProcessEmails()
}
