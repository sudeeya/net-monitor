package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

type snapshotsHTTPServer struct {
	handler *chi.Mux
	service services.SnapshotsService
}

func NewSnapshotsHTTPServer(handler *chi.Mux, service services.SnapshotsService) *snapshotsHTTPServer {
	return &snapshotsHTTPServer{
		handler: handler,
		service: service,
	}
}
