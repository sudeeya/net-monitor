package api

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/converter"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"github.com/sudeeya/net-monitor/internal/server/services"
	"go.uber.org/zap"
)

type snapshotsGRPCServer struct {
	pb.UnimplementedSnapshotsServer
	logger  *zap.Logger
	service services.SnapshotsService
}

func NewSnapshotsGRPCServer(logger *zap.Logger, service services.SnapshotsService) *snapshotsGRPCServer {
	return &snapshotsGRPCServer{
		logger:  logger,
		service: service,
	}
}

func (s *snapshotsGRPCServer) SaveSnapshot(ctx context.Context, request *pb.SaveSnapshotRequest) (*pb.SaveSnapshotResponse, error) {
	var response pb.SaveSnapshotResponse

	if err := s.service.SaveSnapshot(ctx, *converter.ToSnapshotFromProto(request.Snapshot)); err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}
