package apiserver

import "apimod/internal/app/store"

type Config struct {
	BindAddr string `toml:"bind_adddr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Store:    store.NewConfig(),
	}
}
