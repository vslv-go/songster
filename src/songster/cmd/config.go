package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var config Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.WithError(err).Fatal("load dotenv error")
	}

	err = env.Parse(&config)
	if err != nil {
		log.WithError(err).Fatal("init config error")
	}

	initLog(config.LogLevel)
}

type Config struct {
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	PgDSL       string `env:"PG_DSL,required"`
	Port        string `env:"PORT" envDefault:"8080"`
	InfoAPIAddr string `env:"INFO_API_ADDR" envDefault:"localhost:8081"`
}
