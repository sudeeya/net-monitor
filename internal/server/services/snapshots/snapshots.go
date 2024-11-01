package snapshots

import (
	"context"
	"time"

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

func (s *snapshots) DeleteSnapshot(ctx context.Context, id model.ID) error {
	if err := s.repo.DeleteSnapshot(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *snapshots) GetSnapshot(ctx context.Context, id model.ID) (model.Snapshot, error) {
	snapshot, err := s.repo.GetSnapshot(ctx, id)
	if err != nil {
		return model.Snapshot{}, err
	}

	return snapshot, nil
}

func (s *snapshots) GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error) {
	timestamps, err := s.repo.GetNTimestamps(ctx, n)
	if err != nil {
		return nil, err
	}

	return timestamps, nil
}

func (s *snapshots) SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	if err := s.repo.StoreSnapshot(ctx, snapshot); err != nil {
		return err
	}

	return nil
}
