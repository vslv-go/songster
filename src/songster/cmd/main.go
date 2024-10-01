package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	application "songster/app"
	infoClient "songster/info_client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		waitForInterrupt()
		cancel()
	}()

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	app := application.New(db, infoClient.NewInfoClient(config.InfoAPIAddr))

	defer shutdownServer(runServer(cancel, app))

	<-ctx.Done()
}

func waitForInterrupt() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-exit
	log.Debug("app stopped")
}

func runServer(cancel context.CancelFunc, app *application.App) func(ctx context.Context) error {
	server := initRestServer(app)

	go func() {
		defer cancel()
		log.Infof("start rest server at :%s", config.Port)

		if err := server.Run(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Debug("rest server stopped")
				return
			}

			log.WithError(err).Error("rest server error")
		}
	}()

	return server.Shutdown
}

func shutdownServer(stopServerFunc func(ctx context.Context) error) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Debug("stopping rest server")
	if err := stopServerFunc(shutdownCtx); err != nil {
		log.WithError(err).Error("stopping rest server error")
	}
}
