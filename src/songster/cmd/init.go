package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"songster/api/rest"
	"songster/app"
	"songster/repo/pg"
)

const (
	migrationsPath = "repo/pg/migrations"
)

func initLog(logLevel string) {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{})
	if level, err := logrus.ParseLevel(logLevel); err == nil {
		logrus.SetLevel(level)
	} else {
		logrus.WithField("level", logLevel).Warn("invalid log level")
	}
}

func initDB() (*pg.DB, error) {
	db, err := pg.New(config.PgDSL)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	err = db.RunMigrations(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func initRestServer(app *app.App) *rest.Server {
	server := rest.New(app, config.Port)
	server.Init()
	return server
}
