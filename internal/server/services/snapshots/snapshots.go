package snapshots

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"github.com/sudeeya/net-monitor/internal/server/services"
	"go.uber.org/zap"
)

var _ services.SnapshotsService = (*snapshots)(nil)

type snapshots struct {
	logger *zap.Logger
	repo   repository.Repository
}

func NewSnapshots(logger *zap.Logger, repo repository.Repository) *snapshots {
	return &snapshots{
		logger: logger,
		repo:   repo,
	}
}

func (s *snapshots) DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error {
	panic("unimplemented")
}

func (s *snapshots) GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error) {
	panic("unimplemented")
}

func (s *snapshots) ListTimestamps(ctx context.Context) ([]model.Timestamp, error) {
	panic("unimplemented")
}

func (s *snapshots) SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	panic("unimplemented")
}
