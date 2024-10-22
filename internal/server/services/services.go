package services

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type SnapshotsService interface {
	SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error
	GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error)
	ListTimestamps(ctx context.Context) ([]model.Timestamp, error)
	DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error
}
