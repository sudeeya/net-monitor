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
	var snapshotID int
	if err := tx.QueryRow(ctx, insertSnapshotQuery, snapshotArgs).Scan(&snapshotID); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	for _, device := range snapshot.Devices {
		vendorArgs := pgx.NamedArgs{
			"vendor": device.Vendor,
		}
		var vendorID int
		if err := tx.QueryRow(ctx, insertVendorQuery, vendorArgs).Scan(&vendorID); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		osArgs := pgx.NamedArgs{
			"os":      device.OSName,
			"version": device.OSVersion,
		}
		var osID int
		if err := tx.QueryRow(ctx, insertOperatingSystemQuery, osArgs).Scan(&osID); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		deviceArgs := pgx.NamedArgs{
			"vendor_id":           vendorID,
			"operating_system_id": osID,
			"hostname":            device.Hostname,
			"serial_number":       device.Serial,
		}
		var deviceID int
		if err := tx.QueryRow(ctx, insertDeviceQuery, deviceArgs).Scan(&deviceID); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		deviceStateArgs := pgx.NamedArgs{
			"snapshot_id":            snapshotID,
			"device_id":              deviceID,
			"is_snapshot_successful": device.IsSnapshotSuccessful,
		}
		var deviceStateID int
		if err := tx.QueryRow(ctx, insertDeviceStateQuery, deviceStateArgs).Scan(&deviceStateID); err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}

		for _, iface := range device.Interfaces {
			ifaceArgs := pgx.NamedArgs{
				"device_id": deviceID,
				"name":      iface.Name,
			}
			var ifaceID int
			if err := tx.QueryRow(ctx, insertInterfaceQuery, ifaceArgs).Scan(&ifaceID); err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}

			ifaceStateArgs := pgx.NamedArgs{
				"interface_id":    ifaceID,
				"device_state_id": deviceStateID,
				"is_up":           iface.IsUp,
				"ip":              iface.IP,
				"mtu":             iface.MTU,
			}
			if _, err := tx.Exec(ctx, insertInterfaceStateQuery, ifaceStateArgs); err != nil {
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
func (p *postgreSQL) GetNTimestamps(ctx context.Context, n int) ([]model.Snapshot, error) {
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

	timestamps := make([]model.Snapshot, len(dbTimestamps))
	for i, dbt := range dbTimestamps {
		timestamps[i] = model.Snapshot{
			ID:        int(dbt.ID.Int64),
			Timestamp: dbt.Timestamp.Time,
		}
	}

	return timestamps, nil
}

// GetSnapshot implements the [Repository] interface.
func (p *postgreSQL) GetSnapshot(ctx context.Context, id int) (model.Snapshot, error) {
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
func (p *postgreSQL) DeleteSnapshot(ctx context.Context, id int) error {
	p.logger.Info("Deleting a snapshot from the database")

	args := pgx.NamedArgs{
		"id": int(id),
	}
	if _, err := p.db.Exec(ctx, deleteSnapshotQuery, args); err != nil {
		return err
	}

	return nil
}
