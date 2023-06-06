package http

import (
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const salt = "gqgwgd1g21gehwdwh08w7dbb1y2hshsdasd"

type Auth interface {
	SignUp(userAuth entities.UserAuth) (string, error)
	SignIn(userAuth entities.UserAuth) (string, error)
	NewRefreshToken() (string, error)
}

type Adapter struct {
	server *http.Server
	auth   Auth
}

func New(auth Auth) *Adapter {
	return &Adapter{auth: auth}
}

func (a *Adapter) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/signIn", a.signIn)
	mux.HandleFunc("/signUp", a.signUp)

	//chi router почитать

	http.ListenAndServe(
		":8080",
		mux,
	)
	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	return nil
}

func (a *Adapter) signIn(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	password = GeneratePasswordHash(password)
	var userAuth entities.UserAuth
	userAuth.UserName = username
	userAuth.Password = password
	accessToken, err := a.auth.SignIn(userAuth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if accessToken == "Пользователя не существует" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(accessToken))
		return
	}
	refreshToken, err := a.auth.NewRefreshToken()
	response, err := json.Marshal(map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	userAuth.Password = GeneratePasswordHash(userAuth.Password)
	response, err := a.auth.SignUp(userAuth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
