package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"authservice/internal/config"
)

type Redis struct {
	db  *redis.Client
	ttl int
}

func NewClient(config *config.RedisConfig) *Redis {
	rdb := redis.NewClient(
		&redis.Options{
			Addr: fmt.Sprintf("%s:%s", config.Address, config.Port),
		},
	)

	return &Redis{db: rdb, ttl: config.TTL}
}
