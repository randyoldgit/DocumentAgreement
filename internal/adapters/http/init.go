package http

import (
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"time"
)

type Auth interface {
	SignUp(ctx context.Context, userAuth entities.UserAuth) error
	SignIn(ctx context.Context, userAuth entities.UserAuth) (entities.Tokens, error)
	Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error)
	Logout(ctx context.Context, tokens entities.Tokens) error
}

type Adapter struct {
	server *http.Server
	auth   Auth
}

func New(auth Auth) *Adapter {
	return &Adapter{auth: auth}
}

func (a *Adapter) Start() error {
	r := chi.NewRouter()

	r.Use(middleware.Timeout(10 * time.Second))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ping"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/signUp", a.signUp)
		r.Post("/signIn", a.signIn)
		r.Post("/verify", a.verify)
		r.Post("/logout", a.logout)
	})

	http.ListenAndServe(":8080", r)
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	return nil
}

func (a *Adapter) signUp(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	var userAuth entities.UserAuth
	err = json.Unmarshal(body, &userAuth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = a.auth.SignUp(r.Context(), userAuth)
	if errors.Is(err, entities.ErrUserAlreadyExists) || errors.Is(err, entities.ErrInvalidUserCredentials) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	//записать ID пользователя из контекста w.Write([]byte(err.Error()))
}
func (a *Adapter) signIn(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	var userAuth entities.UserAuth
	userAuth.UserName = username
	userAuth.Password = password

	tokens, err := a.auth.SignIn(r.Context(), userAuth)
	if errors.Is(err, entities.ErrUserNotFound) || errors.Is(err, entities.ErrInvalidUserCredentials) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	response, err := json.Marshal(map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (a *Adapter) verify(w http.ResponseWriter, r *http.Request) {
	//TODO Вот эта штука должна идти не в теле а в bearer хедере
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	var tokens entities.Tokens
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokens, err = a.auth.Verify(r.Context(), tokens)
	if errors.Is(err, entities.ErrRefreshTokenInvalid) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (a *Adapter) logout(w http.ResponseWriter, r *http.Request) {
	//TODO Вот эта штука должна идти не в теле а в bearer хедере
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	var tokens entities.Tokens
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
