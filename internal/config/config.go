package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Conf struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	Rabbit  RabbitConf
}

type HTTPConf struct {
	Host string
	Port string
}

type LoggerConf struct {
	Level string
	File  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type RabbitConf struct {
	Dsn      string
	Exchange string
	Queue    string
}

func LoadConfiguration() (*Conf, error) {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	return &Conf{
		Logger: LoggerConf{
			Level: os.Getenv("LOG_LEVEL"),
			File:  os.Getenv("LOG_FILENAME"),
		},
		Storage: StorageConf{
			Type: os.Getenv("STORAGE_TYPE"),
			Dsn:  os.Getenv("STORAGE_DSN"),
		},
		HTTP: HTTPConf{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},
		Rabbit: RabbitConf{
			Dsn:      os.Getenv("RABBIT_DSN"),
			Exchange: os.Getenv("RABBIT_EXCHANGE"),
			Queue:    os.Getenv("RABBIT_QUEUE"),
		},
	}, nil
}
