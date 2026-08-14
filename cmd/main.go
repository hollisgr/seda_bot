package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sedabot/internal/config"
	"sedabot/internal/handler"
	"sedabot/internal/storage/postgres"
	"sedabot/internal/usecase"
	"sedabot/pkg/psqlclient"

	"github.com/go-telegram/bot"
)

func main() {
	log.Println("version 0.0.1a")
	cfg := config.LoadConfig(".env")
	pool, err := psqlclient.NewPool(context.Background(), 3, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalln("postgres new pool err: ", err)
	}

	storage := postgres.New(pool)

	handler := handler.New(
		usecase.NewUserUseCase(storage),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler.Default),
	}

	b, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		panic(err)
	}

	handler.Register(b)
	log.Println("ready to register")

	b.Start(ctx)
}
