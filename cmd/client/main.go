package main

import (
	"log"

	"github.com/sudeeya/net-monitor/internal/client/app"
	"github.com/sudeeya/net-monitor/internal/client/client"
	"github.com/sudeeya/net-monitor/internal/client/config"
	"github.com/sudeeya/net-monitor/internal/client/snapper/snapshots"
	"github.com/sudeeya/net-monitor/internal/pkg/logging"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatal(err)
	}

	snapper, err := snapshots.NewSnapshots(logger, cfg.TargetsFile)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.NewClient(
		cfg.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	snapshotsClient := pb.NewSnapshotsClient(conn)

	grpcClient := client.NewClient(logger, snapper, snapshotsClient)

	a := app.NewApp(cfg, logger, grpcClient)

	a.Run()
}
