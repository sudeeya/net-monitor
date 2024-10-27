package repository

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type Repository interface {
	StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error
	GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error)
	GetNTimestamps(ctx context.Context, n int) ([]model.Timestamp, error)
	DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error
}
