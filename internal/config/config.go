package config

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type PostgresCfg struct {
	Host     string `env:"PSQL_HOST" env-required:"true"`
	Port     string `env:"PSQL_PORT" env-default:"5432"`
	Database string `env:"PSQL_NAME" env-required:"true"`
	Username string `env:"PSQL_USER" env-required:"true"`
	Password string `env:"PSQL_PASSWORD" env-required:"true"`
}

func (p PostgresCfg) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		p.Username, p.Password, p.Host, p.Port, p.Database)
}

type Config struct {
	BotToken string `env:"BOT_TOKEN"`
	OwnerId  int64  `env:"OWNER_ID"`
	Postgres PostgresCfg
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
