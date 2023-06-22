package internal

import (
	repository "DocumentAgreement/internal/adapters/db"
	"DocumentAgreement/internal/adapters/db/postgres"
	"DocumentAgreement/internal/adapters/http"
	"DocumentAgreement/internal/adapters/usecases/auth"
	"context"
	"fmt"
	"log"
)

type (
	App struct {
		server *http.Adapter
	}
)

func (a *App) Start() error {
	pc := postgres.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "123",
		DBName:   "postgres",
		SSLMode:  "disable",
	}
	db, err := postgres.NewPostgresDB(pc)
	if err != nil {
		log.Panic(fmt.Sprintf("Can't connect to DB: %s", err.Error()))
	}
	repo := repository.NewRepository(db)
	authService := auth.New(repo)
	a.server = http.New(authService)
	err = a.server.Start()
	if err != nil {
		//логфатал делаем только в мейне, fmt.errorf <- почитать
		fmt.Println(fmt.Sprintf("Ошибка старта сервера %s", err.Error()))
		log.Fatal(err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.server.Stop(ctx)
	return nil
}
