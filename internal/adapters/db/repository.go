package repository

import (
	"DocumentAgreement/internal/adapters/db/cache_redis"
	"DocumentAgreement/internal/adapters/db/postgres"
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, userAuth entities.UserAuth) error
	GetByCredentials(ctx context.Context, userAuth entities.UserAuth) (int, error)
	IsLoginFree(ctx context.Context, userAuth entities.UserAuth) (bool, error)
	Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error)
}

type AuthorizationCache interface {
	TokenIsValid(ctx context.Context, tokens entities.Tokens) (bool, error)
	SetTokenInvalid(ctx context.Context, tokens entities.Tokens) error
}

type Repository struct {
	Authorization
	AuthorizationCache
}

func NewRepository(db *sqlx.DB, client *redis.Client) *Repository {
	return &Repository{
		Authorization:      postgres.NewAuthPostgres(db),
		AuthorizationCache: cache_redis.NewAuthRedis(client),
	}
}
