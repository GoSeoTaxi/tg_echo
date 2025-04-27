package config

import "github.com/caarlos0/env/v11"

type Config struct {
	BotToken string `env:"BOT_TOKEN,required"`
	ChatID   int64  `env:"CHAT_ID,required"`
	Port     string `env:"PORT" envDefault:"8080"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func Load() *Config {
	var c Config
	if err := env.Parse(&c); err != nil {
		panic(err)
	}
	return &c
}
