package main

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port      int    `env:"PORT" envDefault:"3000"`
	DBDSN     string `env:"DB_DSN"`
	CacherDSN string `env:"CACHER_DSN"`
}

var cfg Config

// InitConfig parses env variables into proper struct
func InitConfig() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		panic("Cannot parse env variables")
	}
}

// GetConfig returns current settings and connection DSNs
func GetConfig() Config {
	return cfg
}
