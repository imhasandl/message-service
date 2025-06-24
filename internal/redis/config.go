package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisConfig(password string) *RedisConfig {
	return &RedisConfig{
		Host:     "localhost",
		Port:     "3663",
		Password: password,
		DB:       0,
	}
}

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedisClient(cfg *RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic("Redis client is not working: " + err.Error())
	}
}
