package api

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/sudeeya/net-monitor/internal/pkg/converter"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

// snapshotsImplementation defines object to interact with the server using gRPC.
type snapshotsImplementation struct {
	pb.UnimplementedSnapshotsServer
	logger  *zap.Logger
	service services.SnapshotsService
}

// NewSnapshotsGRPCServer returns snapshotsGRPCServer object.
func NewSnapshotsGRPCServer(logger *zap.Logger, service services.SnapshotsService) *grpc.Server {
	snapshots := &snapshotsImplementation{
		logger:  logger,
		service: service,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSnapshotsServer(grpcServer, snapshots)

	return grpcServer
}

// SaveSnapshot requsts the service to save the object.
func (s *snapshotsImplementation) SaveSnapshot(ctx context.Context, request *pb.SaveSnapshotRequest) (*pb.SaveSnapshotResponse, error) {
	var response pb.SaveSnapshotResponse

	snapshot, err := converter.ToSnapshotFromProto(request.Snapshot)
	if err != nil {
		return nil, err
	}

	if err := s.service.SaveSnapshot(ctx, *snapshot); err != nil {
		response.Error = err.Error()
		return &response, err
	}

	return &response, nil
}
