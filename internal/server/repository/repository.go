// Package repository defines the interaction with an object storing snapshots.
package repository

import (
	"context"
	"time"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

// Repository describes interaction with an object storing snapshots.
type Repository interface {
	// StoreSnapshot stores a snapshot into Repository.
	// Returns an error if the snapshot could not be stored.
	StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error

	// GetSnapshot returns a snapshot by its id.
	// Returns an error if the snapshot could not be returned.
	GetSnapshot(ctx context.Context, timestamp model.ID) (model.Snapshot, error)

	// GetNTimestamps returns the last n snapshot ids and timestamps.
	// If n is greater than the number of snapshots in the repository, returns all timestamps.
	GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error)

	// DeleteSnapshot deletes a snapshot from Repository by its id.
	DeleteSnapshot(ctx context.Context, timestamp model.ID) error
}
