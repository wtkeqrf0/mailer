package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type (
	Config struct {
		ServiceName string
		Email
		Rabbit
		Mongo
		Logger
	}

	Rabbit struct {
		Consumer        QueueConnection
		LoggerPublisher QueueConnection
		CancelPublisher QueueConnection
		GuzzleLogger    QueueConnection
	}

	Email struct {
		SingleSender EmailConnection
		Mailing      EmailConnection
	}

	Mongo struct {
		URL    string
		DBName string
	}

	Logger struct {
		Level    string
		InFile   bool
		FilePath string
	}

	QueueConnection struct {
		URL       string
		QueueName string
	}

	EmailConnection struct {
		Host       string
		Port       uint16
		Username   string
		Password   string
		ReturnPath string
		Name       string
		PrivateKey []byte `json:"-"`
		ErrorsTo   string
	}
)

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	c := new(Config)
	if err := v.Unmarshal(c); err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
		return nil, err
	}

	if secret, err := os.ReadFile("config/single_sender_private_key.pem"); err != nil {
		log.Printf("single_sender_private_key.pem is not found due %v Dkim is disabled.", err)
	} else {
		c.SingleSender.PrivateKey = secret
	}

	if secret, err := os.ReadFile("config/mailing_private_key.pem"); err != nil {
		log.Printf("mailing_private_key.pem is not found due %v Dkim is disabled.", err)
	} else {
		c.Mailing.PrivateKey = secret
	}

	return c, nil
}
