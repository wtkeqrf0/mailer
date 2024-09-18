package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Email  `yaml:"email"`
		Rabbit `yaml:"rabbit"`
		Mongo  `yaml:"mongo"`
	}

	Rabbit struct {
		Email QueueConnection `yaml:"email"`
	}

	Mongo struct {
		Url    string `yaml:"url"`
		DbName string `yaml:"dbName"`
	}

	QueueConnection struct {
		Url       string `yaml:"url"`
		QueueName string `yaml:"queueName"`
	}

	Email struct {
		Host           string `yaml:"host"`
		Port           uint16 `yaml:"port"`
		Username       string `yaml:"username"`
		Password       string `yaml:"password"`
		ReturnPath     string `yaml:"returnPath"`
		Name           string `yaml:"name"`
		PrivateKeyPath string `yaml:"privateKeyPath"`
		ErrorsTo       string `yaml:"errorsTo"`
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
