package main

import (
	"DocumentAgreement/internal"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//@title DocumentAgreement API
//@ version 1.0
//@description API Server for DocumentAgreement Application

//@host localhost:8080
//@BasePath /

//@securityDefinitions.apikey ApiKeyAuth
//@in header
//@name Authorization

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	app := internal.App{}
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}
