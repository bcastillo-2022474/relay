package config

import "github.com/caarlos0/env/v11"

type Config struct {
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://relay:relay@localhost:5432/relay"`
	HTTPAddr    string `env:"HTTP_ADDR" envDefault:":8080"`
}

func Load() (Config, error) {
	return env.ParseAs[Config]()
}
