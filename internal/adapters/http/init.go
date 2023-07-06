package http

import (
	_ "DocumentAgreement/docs"
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"io"
	"net/http"
	"strings"
	"time"
)

type Auth interface {
	SignUp(ctx context.Context, userAuth entities.UserAuth) error
	SignIn(ctx context.Context, userAuth entities.UserAuth) (entities.Tokens, error)
	Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error)
	Logout(ctx context.Context, tokens entities.Tokens) error
	TokenIsValid(ctx context.Context, tokens entities.Tokens) (bool, error)
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
		r.Post("/refreshToken", a.refreshToken)
	})
	r.Route("/", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Post("/logout", a.logout)
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition
	))

	http.ListenAndServe(":8080", r)
	return nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "Bearer "
		authHeader := r.Header.Get("Authorization")
		reqToken := strings.TrimPrefix(authHeader, prefix)

		if authHeader == "" || reqToken == authHeader {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var tokens entities.Tokens
		tokens.AccessToken = authHeader
		//Пока считаем, что все токены валидны
		//Вот здесь каким-то образом я должен попасть в a.auth.TokenIsValid()
		next.ServeHTTP(w, r)
	})
}

func (a *Adapter) Stop(ctx context.Context) error {
	return nil
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body entities.UserAuth true "list info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /auth/signUp [post]
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

// @Summary SignIn
// @Tags auth
// @Description login in account
// @ID login-account
// @Accept json
// @Produce json
// @Param   username      query     string     false  "string valid"       minlength(5)  maxlength(10)
// @Param   password      query     string     false  "string valid"       minlength(5)  maxlength(10)
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /auth/signIn [post]
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
	response, err := json.MarshalIndent(map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	}, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *Adapter) refreshToken(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(err.Error()))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.MarshalIndent(map[string]interface{}{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	}, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (a *Adapter) logout(w http.ResponseWriter, r *http.Request) {
	var tokens entities.Tokens
	//Рефреш токен
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Записали рефреш токен в редис
	err = a.auth.Logout(r.Context(), tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
