// Package services defines service for interacting with snapshots.
package services

import (
	"context"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

// SnapshotsService describes the service for interacting with snapshots.
type SnapshotsService interface {
	// SaveSnapshot saves a snapshot.
	// Returns an error if the snapshot could not be saved.
	SaveSnapshot(ctx context.Context, snapshot model.Snapshot) error

	// GetSnapshot returns a snapshot by its id.
	// Returns an error if the snapshot could not be returned.
	GetSnapshot(ctx context.Context, id int) (model.Snapshot, error)

	// GetNTimestamps returns the last n snapshot ids and timestamps.
	GetNTimestamps(ctx context.Context, n int) ([]model.Snapshot, error)

	// DeleteSnapshot deletes a snapshot by its id.
	DeleteSnapshot(ctx context.Context, id int) error
}
