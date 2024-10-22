package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

type snapshotsHTTPServer struct {
	*chi.Mux
	service services.SnapshotsService
}

func NewSnapshotsHTTPServer(mux *chi.Mux, service services.SnapshotsService) *snapshotsHTTPServer {
	return &snapshotsHTTPServer{
		Mux:     mux,
		service: service,
	}
}
