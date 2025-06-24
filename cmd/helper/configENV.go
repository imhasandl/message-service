package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	Port        string
	DBURL       string
	TokenSecret string
	RabbitMQ    string
	RedisSecret string
}

func GetENVSecrets() EnvConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}

	config := EnvConfig{
		Port:        os.Getenv("PORT"),
		DBURL:       os.Getenv("DB_URL"),
		TokenSecret: os.Getenv("TOKEN_SECRET"),
		RabbitMQ:    os.Getenv("RABBITMQ_URL"),
		RedisSecret: os.Getenv("REDIS_SECRET"),
	}

	if config.Port == "" {
		log.Fatalf("Set Port in env")
	}
	if config.DBURL == "" {
		log.Fatalf("Set db connection in env")
	}
	if config.TokenSecret == "" {
		log.Fatalf("Set token secret in env")
	}
	if config.RabbitMQ == "" {
		log.Fatalf("Set redis password in .env file")
	}
	if config.RedisSecret == "" {
		log.Fatalf("Set redis password in .env file")
	}

	return config
}
