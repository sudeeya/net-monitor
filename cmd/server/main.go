package main

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/sudeeya/net-monitor/internal/pkg/logging"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"github.com/sudeeya/net-monitor/internal/server/api"
	"github.com/sudeeya/net-monitor/internal/server/app"
	"github.com/sudeeya/net-monitor/internal/server/config"
	"github.com/sudeeya/net-monitor/internal/server/repository/postgresql"
	"github.com/sudeeya/net-monitor/internal/server/services/snapshots"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := postgresql.NewPostgreSQL("")
	if err != nil {
		log.Fatal(err)
	}

	service := snapshots.NewSnapshots(repo)

	grpcServer := grpc.NewServer()
	snapshotsGRPCServer := api.NewSnapshotsGRPCServer(service)
	pb.RegisterSnapshotsServer(grpcServer, snapshotsGRPCServer)

	mux := chi.NewRouter()
	snapshotsHTTPServer := api.NewSnapshotsHTTPServer(mux, service)

	a := app.NewApp(cfg, logger, repo, snapshotsHTTPServer, grpcServer)

	a.Run()
}
