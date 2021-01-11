package main

import (
	"github.com/imarrche/jwt-auth-example/internal/config"
	"github.com/imarrche/jwt-auth-example/internal/logger"
	"github.com/imarrche/jwt-auth-example/internal/server"
	"github.com/imarrche/jwt-auth-example/internal/service/app"
	"github.com/imarrche/jwt-auth-example/internal/store/pg"
)

func main() {
	// Getting logger.
	l := logger.Get()

	// Reading config.
	c := config.Get()

	// Opening PostgreSQL store.
	store := pg.Get(c.PostgreSQL)
	if err := store.Open(); err != nil {
		l.Fatal(err.Error())
	}
	logger.Get().Info("connected to the PostgreSQL")

	// Initalizing service.
	service := app.NewService(store)

	// Running the server.
	server.New(c.Server, service).Run()

	// Closing PostgreSQL store.
	if err := store.Close(); err != nil {
		l.Fatal(err.Error())
	}
}
