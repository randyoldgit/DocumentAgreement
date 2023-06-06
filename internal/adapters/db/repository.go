package repository

import (
	"DocumentAgreement/internal/adapters/entities"
)

type Authorization interface {
	CreateUser(userAuth entities.UserAuth) (string, error)
	SignIn(user entities.UserAuth) (string, error)
}

type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{
		Authorization: NewAuthMemory(),
	}
}
