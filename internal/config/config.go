package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BotToken string `env:"BOT_TOKEN"`
}

func LoadConfig(path string) *Config {
	var cfg Config

	log.Printf("Reading config from %s...", path)

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	log.Println("Config loaded successfully")
	return &cfg
}
