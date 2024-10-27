package services

import (
	"context"
	"time"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type SnapshotsService interface {
	SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error
	GetSnapshot(ctx context.Context, id model.ID) (model.Snapshot, error)
	GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error)
	DeleteSnapshot(ctx context.Context, id model.ID) error
}
