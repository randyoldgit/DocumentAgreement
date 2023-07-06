package internal

import (
	repository "DocumentAgreement/internal/adapters/db"
	"DocumentAgreement/internal/adapters/db/cache_redis"
	"DocumentAgreement/internal/adapters/db/postgres"
	"DocumentAgreement/internal/adapters/entities"
	"DocumentAgreement/internal/adapters/http"
	"DocumentAgreement/internal/adapters/usecases/auth"
	"context"
	"errors"
	"fmt"
	"log"
)

type (
	App struct {
		server *http.Adapter
	}
)

func (a *App) Start() error {
	rc := cache_redis.Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}
	cache, err := cache_redis.NewRedisDB(rc)
	if errors.Is(err, entities.ErrDbConnectionFailed) {
		log.Panic(fmt.Sprintf("Can't connect to Redis: %s", err.Error()))
	}
	if err != nil {
		log.Panic(fmt.Sprintf("Internal Redis error: %s, %s", err.Error()))
	}
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
		log.Panic(fmt.Sprintf("Can't connect to Postgres: %s", err.Error()))
	}
	repo := repository.NewRepository(db, cache)
	authService := auth.New(repo)
	a.server = http.New(authService)
	err = a.server.Start()
	if err != nil {
		fmt.Println(fmt.Sprintf("Ошибка старта сервера %s", err.Error()))
		log.Fatal(err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	a.server.Stop(ctx)
	return nil
}
