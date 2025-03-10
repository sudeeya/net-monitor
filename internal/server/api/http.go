package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/server/handlers"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

// Endpoints.
const (
	defaultEndpoint       = "/"
	getTimestampsEndpoint = "/timestamps"
	getSnapshotEndpoint   = "/snapshots"
)

// snapshotsHTTPServer defines object to interact with the server using HTTP.
type snapshotsHTTPServer struct {
	*chi.Mux
	logger  *zap.Logger
	service services.SnapshotsService
}

// NewSnapshotsHTTPServer returns snapshotsHTTPServer object.
func NewSnapshotsHTTPServer(logger *zap.Logger, service services.SnapshotsService) *snapshotsHTTPServer {
	mux := chi.NewRouter()

	registerEndpoints(mux, logger, service)

	return &snapshotsHTTPServer{
		Mux:     mux,
		logger:  logger,
		service: service,
	}
}

// registerEndpoints registers enpoints for HTTP requests.
func registerEndpoints(mux *chi.Mux, logger *zap.Logger, service services.SnapshotsService) {
	mux.Get(defaultEndpoint, handlers.DefaultHandler(logger))
	mux.Get(getTimestampsEndpoint, handlers.GetTimestampsHandler(logger, service))
	mux.Get(getSnapshotEndpoint, handlers.GetSnapshotHandler(logger, service))
}
