package app

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/server/config"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"google.golang.org/grpc"
)

type app struct {
	cfg        *config.Config
	logger     *zap.Logger
	repo       repository.Repository
	handler    http.Handler
	grpcServer *grpc.Server
}

func NewApp(
	cfg *config.Config,
	logger *zap.Logger,
	repo repository.Repository,
	handler http.Handler,
	grpcServer *grpc.Server,
) *app {
	return &app{
		cfg:        cfg,
		logger:     logger,
		repo:       repo,
		handler:    handler,
		grpcServer: grpcServer,
	}
}

func (s *app) Run() {
	panic("unimplemented")
}
