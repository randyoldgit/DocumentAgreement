package auth

import (
	repository "DocumentAgreement/internal/adapters/db"
	"DocumentAgreement/internal/adapters/entities"
	"DocumentAgreement/internal/hasher"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const (
	signingKey      = "hhqwghd12ejjsad7axucjn12eh1jia8dyas"
	accessTokenTTL  = 5 * time.Minute
	refreshTokenTTL = 60 * time.Minute
)

type Authorization interface {
	CreateUser(ctx context.Context, userAuth entities.UserAuth) error
	GetByCredentials(ctx context.Context, userAuth entities.UserAuth) (int, error)
	IsLoginFree(ctx context.Context, userAuth entities.UserAuth) (bool, error)
	Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error)
	TokenIsValid(ctx context.Context, tokens entities.Tokens) (bool, error)
	SetTokenInvalid(ctx context.Context, tokens entities.Tokens) error
}

type Service struct {
	Authorization
}

func newClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"exp":    0,
		"iat":    0,
		"userId": 0,
	}
}

func New(repo *repository.Repository) *Service {
	return &Service{Authorization: repo}
}

func (s *Service) SignUp(ctx context.Context, userAuth entities.UserAuth) error {
	if len(userAuth.UserName) < 4 {
		return fmt.Errorf("Username must contain at least 4 symbols: %w", entities.ErrInvalidUserCredentials)
	}
	if len(userAuth.Password) < 4 {
		return fmt.Errorf("Password must contain at least 4 symbols: %w", entities.ErrInvalidUserCredentials)
	}
	exists, err := s.IsLoginFree(ctx, userAuth)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w", entities.ErrUserAlreadyExists)
	}
	userAuth.Password = hasher.GeneratePasswordHash(userAuth.Password)
	err = s.CreateUser(ctx, userAuth)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) SignIn(ctx context.Context, userAuth entities.UserAuth) (entities.Tokens, error) {
	var tokens entities.Tokens
	if len(userAuth.UserName) == 0 {
		return tokens, fmt.Errorf("Empty username: %w", entities.ErrInvalidUserCredentials)
	}
	if len(userAuth.Password) == 0 {
		return tokens, fmt.Errorf("Empty password: %w", entities.ErrInvalidUserCredentials)
	}
	userAuth.Password = hasher.GeneratePasswordHash(userAuth.Password)
	userId, err := s.Authorization.GetByCredentials(ctx, userAuth)
	if err != nil {
		return tokens, err
	}
	if userId == 0 {
		return tokens, fmt.Errorf("User doesn't exists: %w", entities.ErrUserNotFound)
	}
	tokens, err = s.createSession(userId)
	if err != nil {
		return tokens, err
	}
	return tokens, nil
}
func (s *Service) Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error) {
	var newTokens entities.Tokens
	//Проверяем, что аксесс валидный, если да - возвращаем его же 0 ошибок
	accessIsValid, err := verifyToken(tokens.AccessToken)
	if accessIsValid {
		return tokens, nil
	}
	//Если не валидный - проверяем рефреш токен
	refreshIsValid, err := verifyToken(tokens.RefreshToken)
	//Если рефреш невалидный - то возвращаем свою ошибку
	if !refreshIsValid {
		return newTokens, fmt.Errorf("%w", entities.ErrRefreshTokenInvalid)
	}
	//Если по другой причине - то просто возвращаем ошибку
	if err != nil {
		return newTokens, err
	}
	//Если рефреш валиден - то генерируем новую пару токенов
	//Достаем userId из рефреш токена
	payload := newClaims()
	_, err = jwt.ParseWithClaims(tokens.RefreshToken, payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	userId, err := strconv.Atoi(fmt.Sprintf("%v", payload["userId"]))
	if err != nil {
		return newTokens, err
	}
	newTokens, err = s.createSession(userId)
	return newTokens, nil
}

func (s *Service) TokenIsValid(ctx context.Context, tokens entities.Tokens) (bool, error) {
	result, err := s.Authorization.TokenIsValid(ctx, tokens)
	if err != nil {
		return false, fmt.Errorf("%w", err)
	}
	return result, nil
}

func (s *Service) Logout(ctx context.Context, tokens entities.Tokens) error {
	err := s.Authorization.SetTokenInvalid(ctx, tokens)
	if err != nil {
		return fmt.Errorf("%w", err.Error())
	}
	return nil
}

func (s *Service) createSession(userId int) (entities.Tokens, error) {
	var res entities.Tokens
	var err error

	res.AccessToken, err = s.newToken(userId, accessTokenTTL)
	if err != nil {
		return res, err
	}
	res.RefreshToken, err = s.newToken(userId, refreshTokenTTL)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *Service) newToken(userId int, d time.Duration) (string, error) {
	payload := newClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	res, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return res, nil
}

func verifyToken(input string) (bool, error) {
	payload := newClaims()
	token, err := jwt.ParseWithClaims(input, payload, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}
