package app

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

func (a *app) Run() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		if err := http.ListenAndServe(a.cfg.HTTPAddr, a.handler); err != nil {
			a.logger.Fatal(err.Error())
		}
	}()

	go func() {
		listen, err := net.Listen("tcp", a.cfg.GRPCAddr)
		if err != nil {
			a.logger.Fatal(err.Error())
		}

		if err := a.grpcServer.Serve(listen); err != nil {
			a.logger.Fatal(err.Error())
		}
	}()

	<-sigCh
	a.Shutdown()
}

func (a *app) Shutdown() {
	if err := a.logger.Sync(); err != nil {
		a.logger.Sugar().Fatalf("failed to sync logger: %v", err)
	}

	os.Exit(0)
}
