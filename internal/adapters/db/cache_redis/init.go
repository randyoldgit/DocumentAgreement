package cache_redis

import (
	"DocumentAgreement/internal/adapters/entities"
	"fmt"
	"github.com/go-redis/redis"
)

type Config struct {
	Host     string
	Port     string
	DB       int
	Password string
}

func NewRedisDB(cfg Config) (*redis.Client, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	pong, err := rc.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s:", entities.ErrDbConnectionFailed, err, pong)
	}
	return rc, nil
}
