package repository

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type Repository interface {
	GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error)
	StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error
	DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error
}
