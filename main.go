package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	cfg "github.com/erazr/go_tz/config"
	db "github.com/erazr/go_tz/db"
	_ "github.com/erazr/go_tz/docs"
	api "github.com/erazr/go_tz/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")

	ctx := context.Background()
	postgresURL, err := cfg.LoadDBConfig()

	m, err := migrate.New("file://db/migrations", postgresURL+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	database, err := db.NewDatabase(ctx, postgresURL)

	api.RunHttp(ctx, database)

	if err != nil {
		panic(err)
	}

	<-stop
	log.Println("Shutting down server...")
	ctx.Done()
}
