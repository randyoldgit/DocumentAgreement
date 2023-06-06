package internal

import (
	repository "DocumentAgreement/internal/adapters/db"
	"DocumentAgreement/internal/adapters/http"
	"DocumentAgreement/internal/adapters/usecases/auth"
	"context"
	"log"
)

type (
	App struct {
		server *http.Adapter
	}
)

func (a *App) Start() error {
	memoryRepo := repository.NewRepository()
	authService := auth.New(memoryRepo)
	a.server = http.New(authService)
	err := a.server.Start()
	if err != nil {
		//логфатал делаем только в мейне, fmt.errorf <- почитать
		log.Fatal(err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.server.Stop(ctx)
	return nil
}
