package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pullya/unique_server/u-server/internal/app"
	"github.com/pullya/unique_server/u-server/internal/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(config.LogLevel)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.WithField("service", config.ServiceName).Warnf("Received signal %v. Shutting down...", sig)

		cancel()
	}()

	app := app.InitApp()

	if err := app.Run(ctx); err != nil {
		log.WithField("service", config.ServiceName).Fatalf("Error while running application: %v", err)
	}
}
