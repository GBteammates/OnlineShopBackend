package jwtauth

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type JWTKey struct {
	Key string `json:"jwtKey" env:"JWTKEY" envDefault:"dsf498uh324seyu2837912sd7*7897"`
}

func NewJWTKeyConfig() (*JWTKey, error) {
	var cfg = JWTKey{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("can't load environment variables: %s", err)
	}
	return &cfg, nil
}
