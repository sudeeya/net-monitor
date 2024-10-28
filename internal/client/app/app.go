package app

import (
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
	go func() {
		for {
			err := a.client.UploadSnapshot()
			if err != nil {
				a.logger.Error(err.Error())
			}
		}
	}()
}
