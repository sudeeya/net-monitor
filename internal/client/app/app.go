package app

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sudeeya/net-monitor/internal/client/client"
	"github.com/sudeeya/net-monitor/internal/client/config"
	"go.uber.org/zap"
)

type app struct {
	cfg    *config.Config
	logger *zap.Logger
	client *client.Client
}

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

func (a *app) Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		for {
			err := a.client.UploadSnapshot()
			if err != nil {
				a.logger.Error(err.Error())
			}
		}
	}()

	<-sigCh
	a.Shutdown()
}

func (a *app) Shutdown() {
	var pathErr fs.PathError
	if err := a.logger.Sync(); err != nil && errors.Is(err, &pathErr) {
		log.Fatalf("failed to sync logger: %v\n", err)
	}

	os.Exit(0)
}
