// Package snapshots defines service that interacts with a [Repository] object.
package snapshots

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

var _ services.SnapshotsService = (*snapshots)(nil)

// snapshots implements the [SnapshotsService] interface.
type snapshots struct {
	logger *zap.Logger
	repo   repository.Repository
}

// NewSnapshots returns snapshots object to interact with a [Repository] object.
func NewSnapshots(logger *zap.Logger, repo repository.Repository) *snapshots {
	return &snapshots{
		logger: logger,
		repo:   repo,
	}
}

// DeleteSnapshot implements the [SnapshotsService] interface.
func (s *snapshots) DeleteSnapshot(ctx context.Context, id model.ID) error {
	s.logger.Info("Deleting a snapshot")
	if err := s.repo.DeleteSnapshot(ctx, id); err != nil {
		return err
	}

	return nil
}

// GetSnapshot implements the [SnapshotsService] interface.
func (s *snapshots) GetSnapshot(ctx context.Context, id model.ID) (model.Snapshot, error) {
	s.logger.Info("Getting a snapshot")
	snapshot, err := s.repo.GetSnapshot(ctx, id)
	if err != nil {
		return model.Snapshot{}, err
	}

	return snapshot, nil
}

// GetNTimestamps implements the [SnapshotsService] interface.
func (s *snapshots) GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error) {
	s.logger.Sugar().Infof("Getting the last %d timestamps", n)
	timestamps, err := s.repo.GetNTimestamps(ctx, n)
	if err != nil {
		return nil, err
	}

	return timestamps, nil
}

// SaveSnapshot implements the [SnapshotsService] interface.
func (s *snapshots) SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	s.logger.Info("Saving a snapshot")
	if err := s.repo.StoreSnapshot(ctx, snapshot); err != nil {
		return err
	}

	return nil
}
