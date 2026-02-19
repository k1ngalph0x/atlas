package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB PostgresConfig
	TOKEN TokenConfig
	KAFKA KafkaConfig
	RABBITMQ RabbitConfig
	OLLAMA OllamaConfig
}

type OllamaConfig struct{
	Url string
	Model string
}

type RabbitConfig struct{
	User string
	Password string
	Host string
	Port string
}

type KafkaConfig struct{
	Brokers []string
}

type TokenConfig struct{
	JwtKey string
}

type PostgresConfig struct {
	Host     string
	Dbname   string
	Username string
	Password string
	Url      string
	Port     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil{
		return nil, err
	}

	config := &Config{
		DB: PostgresConfig{
			Host: os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Url: os.Getenv("DB_URL"),
			Port: os.Getenv("DB_PORT"),
			Dbname: os.Getenv("DB_NAME"),
		},

		TOKEN: TokenConfig{
			JwtKey: os.Getenv("JwtKey"),
		},

		KAFKA: KafkaConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		},

		RABBITMQ: RabbitConfig{
			User: os.Getenv("RABBIT_USER"),
			Password: os.Getenv("RABBIT_PASSWORD"),
			Host: os.Getenv("RABBIT_HOST"),
			Port: os.Getenv("RABBIT_PORT"),
		},
		OLLAMA: OllamaConfig{
			Url: os.Getenv("OLLAMA_URL"),
			Model: os.Getenv("OLLAMA_MODEL"),
		},
	}

	return config, nil

}