package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/server/services"
	"go.uber.org/zap"
)

type snapshotsHTTPServer struct {
	*chi.Mux
	logger  *zap.Logger
	service services.SnapshotsService
}

func NewSnapshotsHTTPServer(logger *zap.Logger, service services.SnapshotsService) *snapshotsHTTPServer {
	mux := chi.NewRouter()

	return &snapshotsHTTPServer{
		Mux:     mux,
		logger:  logger,
		service: service,
	}
}
