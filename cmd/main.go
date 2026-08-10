package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sedabot/internal/config"
	"sedabot/internal/handler"
	"sedabot/pkg/postgres"

	"github.com/go-telegram/bot"
)

func main() {
	cfg := config.LoadConfig(".env")
	pool, err := postgres.NewPool(context.Background(), 3, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalln("postgres new pool err: ", err)
	}
	handler := handler.New()

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

	b.Start(ctx)
}
