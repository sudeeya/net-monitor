package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
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
	envFile := flag.String("e", "env/server.env", "Path to the file storing environment variables")

	flag.Parse()

	if err := godotenv.Load(*envFile); err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := postgresql.NewPostgreSQL(logger, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	service := snapshots.NewSnapshots(logger, repo)

	grpcServer := grpc.NewServer()
	snapshotsGRPCServer := api.NewSnapshotsGRPCServer(logger, service)
	pb.RegisterSnapshotsServer(grpcServer, snapshotsGRPCServer)

	httpServer := api.NewSnapshotsHTTPServer(logger, service)

	a := app.NewApp(cfg, logger, repo, httpServer, grpcServer)

	a.Run()
}
