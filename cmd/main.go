package main

import (
	"context"
	"os"
	"os/signal"
	"sedabot/internal/config"
	"sedabot/internal/handler"

	"github.com/go-telegram/bot"
)

func main() {
	cfg := config.LoadConfig(".env")
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
