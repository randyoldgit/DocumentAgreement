package repository

import (
	"DocumentAgreement/internal/adapters/db/postgres"
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, userAuth entities.UserAuth) error
	GetByCredentials(ctx context.Context, userAuth entities.UserAuth) (int, error)
	IsLoginFree(ctx context.Context, userAuth entities.UserAuth) (bool, error)
	Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: postgres.NewAuthPostgres(db),
	}
}
