package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Server
		Email
		Rabbit
		Mongo
		Logger
	}

	Server struct {
		Name string
	}

	Rabbit struct {
		Mails QueueConnection
		Logs  QueueConnection
	}

	Email struct {
		SingleSender EmailConnection
		Mailing      EmailConnection
	}

	Mongo struct {
		URL    string
		DbName string
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
		Host           string
		Port           uint16
		Username       string
		Password       string
		ReturnPath     string
		Name           string
		PrivateKeyPath string
		ErrorsTo       string
	}
)

func ReadConfigFromFile(configFilePath string) *Config {
	file, err := os.Open(configFilePath)
	if err != nil {
		panic(err)
	}

	cfg := new(Config)
	if err = yaml.NewDecoder(file).Decode(cfg); err != nil {
		panic(err)
	}
	return cfg
}
