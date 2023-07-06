package entities

import "errors"

var (
	ErrDbConnectionFailed     = errors.New("Connection to DB failed")
	ErrUserNotFound           = errors.New("User doesn't exists")
	ErrUserAlreadyExists      = errors.New("User with such username already exists")
	ErrInvalidUserCredentials = errors.New("Username or password isn't correct")
	ErrAccessTokenInvalid     = errors.New("Access token is invalid")
	ErrRefreshTokenInvalid    = errors.New("Refresh token is invalid")
)
