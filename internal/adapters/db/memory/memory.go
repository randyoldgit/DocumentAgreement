package memory

import (
	"DocumentAgreement/internal/adapters/entities"
	"DocumentAgreement/internal/hasher"
	"context"
	"sync"
	"time"
)

type AuthMemory struct {
	users map[string]string
	rwm   sync.RWMutex
}

func NewAuthMemory() *AuthMemory {
	users := make(map[string]string)
	users["Alexander"] = hasher.GeneratePasswordHash("123")
	return &AuthMemory{users: users}
}

func (a *AuthMemory) CreateUser(ctx context.Context, userAuth entities.UserAuth) (int, error) {
	a.rwm.RLock()
	_, exists := a.users[userAuth.UserName]
	a.rwm.RLock()
	//здесь должна ставиться блокировка на запись
	if exists {
		return 0, nil
	}
	a.users[userAuth.UserName] = userAuth.Password
	return 1, nil
}

func (a *AuthMemory) GetUser(user entities.UserAuth) (string, error) {
	//здесь должна ставиться блокировка на чтение, прочитать про мьютексы (rw mutex)
	pass := a.users[user.UserName]
	if user.Password != pass {
		return "Пользователя не существует", nil
	}
	return "Залогинились", nil
}

func (a *AuthMemory) UpdateSession(refreshToken string, expireAt time.Time) (bool, error) {
	return false, nil
}
func (a *AuthMemory) GetSession(refreshToken string) (time.Time, error) {
	return time.Now(), nil
}
