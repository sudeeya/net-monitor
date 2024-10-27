package repository

import (
	"context"
	"time"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type Repository interface {
	StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error
	GetSnapshot(ctx context.Context, timestamp model.ID) (model.Snapshot, error)
	GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error)
	DeleteSnapshot(ctx context.Context, timestamp model.ID) error
}
