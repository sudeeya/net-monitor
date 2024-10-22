package api

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/converter"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

type snapshotsGRPCServer struct {
	pb.UnimplementedSnapshotsServer
	service services.SnapshotsService
}

func (s *snapshotsGRPCServer) SaveSnapshot(ctx context.Context, request *pb.SaveSnapshotRequest) (*pb.SaveSnapshotResponse, error) {
	var response pb.SaveSnapshotResponse

	if err := s.service.SaveSnapshot(ctx, *converter.ToSnapshotFromProto(request.Snapshot)); err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}

func NewSnapshotsGRPCServer(service services.SnapshotsService) *snapshotsGRPCServer {
	return &snapshotsGRPCServer{
		service: service,
	}
}
