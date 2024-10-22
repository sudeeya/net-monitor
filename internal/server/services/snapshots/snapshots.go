package snapshots

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"github.com/sudeeya/net-monitor/internal/server/services"
)

var _ services.SnapshotsService = (*Snapshots)(nil)

type Snapshots struct {
	repo repository.Repository
}

func (s *Snapshots) DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error {
	panic("unimplemented")
}

func (s *Snapshots) GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error) {
	panic("unimplemented")
}

func (s *Snapshots) ListTimestamps(ctx context.Context) ([]model.Timestamp, error) {
	panic("unimplemented")
}

func (s *Snapshots) SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	panic("unimplemented")
}
