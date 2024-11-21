// Package app defines client application object.
package app

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/client/client"
	"github.com/sudeeya/net-monitor/internal/client/config"
)

// app describes client application and all necessary layers.
type app struct {
	cfg    *config.Config
	logger *zap.Logger
	client *client.Client
}

// NewApp returns app object to interact with client.
func NewApp(
	cfg *config.Config,
	logger *zap.Logger,
	client *client.Client,
) *app {
	return &app{
		cfg:    cfg,
		logger: logger,
		client: client,
	}
}

// Run starts the client.
// It initiates periodic sending of gRPC requests to server and monitors for OS signals.
func (a *app) Run() {
	a.logger.Info("Client is running")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	uploadTicker := time.NewTicker(a.cfg.SnapInterval)

	go func() {
		for range uploadTicker.C {
			a.logger.Info("Client is getting ready to upload a snapshot")
			err := a.client.UploadSnapshot()
			if err != nil {
				a.logger.Error(err.Error())
			}
			a.logger.Info("Client uploaded a snapshot")
		}
	}()

	<-sigCh
	a.logger.Info("Client is shutting down")
	a.Shutdown()
}

// Shutdown shuts down the client.
// It syncs client logger before shutdown.
func (a *app) Shutdown() {
	if err := a.client.Close(); err != nil {
		log.Printf("Failed to close the client: %v/n", err)
	}

	var pathErr fs.PathError
	if err := a.logger.Sync(); err != nil && errors.Is(err, &pathErr) {
		log.Printf("Failed to sync logger: %v\n", err)
	}

	os.Exit(0)
}
