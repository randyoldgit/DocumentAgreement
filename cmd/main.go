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
