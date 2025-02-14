// Package postgresql defines object that stores snapshots in PostgeSQL database.
package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
)

const limitInSeconds = 5

var _ repository.Repository = (*postgreSQL)(nil)

// postgreSQL implements the [Repository] interface.
type postgreSQL struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

// NewPostgreSQL returns postgreSQL object to interact with PostgreSQL database.
// The function establishes and tests the connection to the database and creates the necessary tables.
func NewPostgreSQL(logger *zap.Logger, dsn string) (*postgreSQL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
	defer cancel()

	logger.Info("Establishing a connection to the database")
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	logger.Info("Creating necessary tables")
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	createTableQueries := []string{
		createTableSnapshotsQuery,
		createTableVendorsQuery,
		createTableOperatingSystemsQuery,
		createTableDevicesQuery,
		createTableDeviceStatesQuery,
		createTableInterfacesQuery,
		createTableInterfaceStatesQuery,
		createTableSubinterfacesQuery,
		createTableSubinterfaceStatesQuery,
	}

	for _, query := range createTableQueries {
		if _, err := tx.Exec(ctx, query); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return nil, rollbackErr
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &postgreSQL{
		logger: logger,
		db:     db,
	}, nil
}

// StoreSnapshot implements the [Repository] interface.
func (p *postgreSQL) StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	p.logger.Info("Storing a snapshot to the database")

	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	snapshotArgs := pgx.NamedArgs{
		"timestamp": snapshot.Timestamp,
	}
	if _, err := tx.Exec(ctx, insertSnapshotQuery, snapshotArgs); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	for _, device := range snapshot.Devices {
		vendorArgs := pgx.NamedArgs{
			"vendor": device.Vendor,
		}
		if _, err := tx.Exec(ctx, insertVendorQuery, vendorArgs); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		deviceArgs := pgx.NamedArgs{
			"timestamp": snapshot.Timestamp,
			"vendor":    device.Vendor,
			"hostname":  device.Hostname,
			"os":        device.OSName,
			"version":   device.OSVersion,
			"serial":    device.Serial,
			"ip":        device.ManagementIP,
		}
		if _, err := tx.Exec(ctx, insertDeviceQuery, deviceArgs); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		for _, iface := range device.Interfaces {
			ifaceArgs := pgx.NamedArgs{
				"name":      iface.Name,
				"mac":       iface.MAC,
				"ip":        iface.IP,
				"mtu":       iface.MTU,
				"bandwidth": iface.Bandwidth,
			}
			if _, err := tx.Exec(ctx, insertInterfaceQuery, ifaceArgs); err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}
		}
	}

	return tx.Commit(ctx)
}

// GetNTimestamps implements the [Repository] interface.
func (p *postgreSQL) GetNTimestamps(ctx context.Context, n int) (map[model.ID]time.Time, error) {
	p.logger.Sugar().Infof("Getting the last %d timestamps from the database", n)

	args := pgx.NamedArgs{
		"limit": n,
	}
	rows, err := p.db.Query(ctx, selectTimestampsQuery, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbTimestamps, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbTimestamp])
	if err != nil {
		return nil, err
	}

	timestamps := make(map[model.ID]time.Time, len(dbTimestamps))
	for _, dbt := range dbTimestamps {
		timestamps[model.ID(dbt.ID)] = dbt.Timestamp
	}

	return timestamps, nil
}

// GetSnapshot implements the [Repository] interface.
func (p *postgreSQL) GetSnapshot(ctx context.Context, id model.ID) (model.Snapshot, error) {
	p.logger.Info("Getting a snapshot from the database")

	args := pgx.NamedArgs{
		"id": int(id),
	}
	rows, err := p.db.Query(ctx, selectSnapshotQuery, args)
	if err != nil {
		return model.Snapshot{}, err
	}
	defer rows.Close()

	dbSnapshotParts, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbSnapshotPart])
	if err != nil {
		return model.Snapshot{}, err
	}

	return toSnapshotFromDB(dbSnapshotParts), nil
}

// DeleteSnapshot implements the [Repository] interface.
func (p *postgreSQL) DeleteSnapshot(ctx context.Context, id model.ID) error {
	p.logger.Info("Deleting a snapshot from the database")

	args := pgx.NamedArgs{
		"id": int(id),
	}
	if _, err := p.db.Exec(ctx, deleteSnapshotQuery, args); err != nil {
		return err
	}

	return nil
}
