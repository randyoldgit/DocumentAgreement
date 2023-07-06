package cache_redis

import (
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type AuthRedis struct {
	client *redis.Client
}

const refreshTokenTTL = 60 * time.Minute

func NewAuthRedis(client *redis.Client) *AuthRedis {
	return &AuthRedis{client: client}
}

func (r *AuthRedis) TokenIsValid(ctx context.Context, tokens entities.Tokens) (bool, error) {
	val, err := r.client.Get(tokens.RefreshToken).Result()
	if val == "" {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return true, nil
}
func (r *AuthRedis) SetTokenInvalid(ctx context.Context, tokens entities.Tokens) error {
	err := r.client.Set("expiredToken", tokens.RefreshToken, refreshTokenTTL).Err()
	fmt.Println(err)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
