// Package services defines service for interacting with snapshots.
package services

import (
	"context"
	"time"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

// SnapshotsService describes the service for interacting with snapshots.
type SnapshotsService interface {
	// SaveSnapshot saves a snapshot.
	// Returns an error if the snapshot could not be saved.
	SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error

	// GetSnapshot returns a snapshot by its id.
	// Returns an error if the snapshot could not be returned.
	GetSnapshot(ctx context.Context, id model.ID) (model.Snapshot, error)

	// GetNTimestamps returns the last n snapshot ids and timestamps.
	GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error)

	// DeleteSnapshot deletes a snapshot by its id.
	DeleteSnapshot(ctx context.Context, id model.ID) error
}
