package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sudeeya/net-monitor/internal/client/app"
	"github.com/sudeeya/net-monitor/internal/client/client"
	"github.com/sudeeya/net-monitor/internal/client/config"
	"github.com/sudeeya/net-monitor/internal/client/snapper/snapshots"
	"github.com/sudeeya/net-monitor/internal/pkg/logging"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
)

func main() {
	_, err := fmt.Printf("Build version: %s\nBuild date: %s\n", buildVersion, buildDate)
	if err != nil {
		log.Fatal(err)
	}

	envFile := flag.String("e", "env/client.env", "Path to the file storing environment variables")

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

	snapper, err := snapshots.NewSnapshots(logger, cfg.TargetsFile)
	if err != nil {
		log.Fatal(err)
	}

	grpcClient, err := client.NewClient(logger, snapper, cfg.ServerAddr)
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(cfg, logger, grpcClient)

	a.Run()
}
