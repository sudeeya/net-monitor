package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/server/handlers"
	"github.com/sudeeya/net-monitor/internal/server/services"
	"go.uber.org/zap"
)

const (
	getNTimestampsEndpoint = "/timestamps/{timestampsCount}"
	getSnapshotEndpoint    = "/snapshot/{snapshotID}"
	deleteSnapshotEndpoint = "/delete/{snapshotID}"
)

type snapshotsHTTPServer struct {
	*chi.Mux
	logger  *zap.Logger
	service services.SnapshotsService
}

func NewSnapshotsHTTPServer(logger *zap.Logger, service services.SnapshotsService) *snapshotsHTTPServer {
	mux := chi.NewRouter()

	registerEndpoints(mux, logger, service)

	return &snapshotsHTTPServer{
		Mux:     mux,
		logger:  logger,
		service: service,
	}
}

func registerEndpoints(mux *chi.Mux, logger *zap.Logger, service services.SnapshotsService) {
	mux.Get(getNTimestampsEndpoint, handlers.GetNTimestampsHandler(logger, service))
	mux.Get(getSnapshotEndpoint, handlers.GetSnapshotHandler(logger, service))
	mux.Get(deleteSnapshotEndpoint, handlers.DeleteSnapshotHandler(logger, service))
}
