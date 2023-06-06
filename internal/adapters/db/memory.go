package repository

import (
	"DocumentAgreement/internal/adapters/entities"
	"DocumentAgreement/internal/adapters/http"
)

type AuthMemory struct {
	users map[string]string
}

func NewAuthMemory() *AuthMemory {
	users := make(map[string]string)
	users["Alexander"] = http.GeneratePasswordHash("123")
	return &AuthMemory{users: users}
}

func (a *AuthMemory) CreateUser(user entities.UserAuth) (string, error) {
	_, exists := a.users[user.UserName]
	//здесь должна ставиться блокировка на запись
	if exists {
		return "Логин занят", nil
	}
	a.users[user.UserName] = user.Password
	return "Пользователь был создан", nil
}

func (a *AuthMemory) SignIn(user entities.UserAuth) (string, error) {
	//здесь должна ставиться блокировка на чтение, прочитать про мьютексы (rw mutex)
	pass := a.users[user.UserName]
	if user.Password != pass {
		return "Пользователя не существует", nil
	}
	return "Залогинились", nil
}
