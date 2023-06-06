package auth

import (
	repository "DocumentAgreement/internal/adapters/db"
	"DocumentAgreement/internal/adapters/entities"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

const (
	signingKey     = "hhqwghd12ejjsad7axucjn12eh1jia8dyas"
	accessTokenTTL = 5 * time.Minute
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignUp(userAuth entities.UserAuth) (string, error) {
	status, err := s.repo.CreateUser(userAuth)
	if err != nil {
		return status, err
	}
	return status, nil
}

type tokenClaims struct {
	jwt.StandardClaims
	UserName string `json:"login"`
}

func (s *Service) SignIn(userAuth entities.UserAuth) (string, error) {
	result, err := s.repo.SignIn(userAuth)
	if err != nil {
		return "", err
	}
	if result == "Пользователя не существует" {
		return result, nil
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(accessTokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userAuth.UserName,
	})
	return accessToken.SignedString([]byte(signingKey))
}

func (s *Service) NewRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	source := rand.NewSource(time.Now().Unix())
	rand := rand.New(source)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", bytes), nil
}
