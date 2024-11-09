// Package postgresql defines object that stores snapshots in PostgeSQL database.
package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"go.uber.org/zap"
)

const limitInSeconds = 5

// SQL queries to create tables needed to store snapshots.
const (
	createTableSnapshotsQuery = `
CREATE TABLE IF NOT EXISTS snapshots (
	id SERIAL PRIMARY KEY,
	timestamp TIMESTAMP UNIQUE NOT NULL
);
`

	createTableVendorsQuery = `
CREATE TABLE IF NOT EXISTS vendors (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL
);
`

	createTableDevicesQuery = `
CREATE TABLE IF NOT EXISTS devices (
	id SERIAL PRIMARY KEY,
	snapshot_id INT REFERENCES snapshots(id) ON DELETE CASCADE,
	vendor_id INT REFERENCES vendors(id) ON DELETE RESTRICT,
	hostname varchar(255),
	os_name VARCHAR(255),
	os_version VARCHAR(255),
	serial_number VARCHAR(255),
	management_ip INET
);
`

	createTableInterfacesQuery = `
CREATE TABLE IF NOT EXISTS interfaces (
	id SERIAL PRIMARY KEY,
	device_id INT REFERENCES devices(id) ON DELETE CASCADE,
	name VARCHAR(255),
	mac MACADDR,
	ip INET,
	mtu INT,
	bandwidth INT
);
`
)

// SQL queries for inserting a snapshot.
const (
	insertSnapshotQuery = `
INSERT INTO snapshots (timestamp)
VALUES (@timestamp);
`

	insertVendorQuery = `
INSERT INTO vendors (name)
VALUES (@vendor)
ON CONFLICT (name) DO NOTHING;
`

	insertDeviceQuery = `
WITH snapshot_id AS (
	SELECT id
	FROM snapshots
	WHERE timestamp = @timestamp
), vendor_id AS (
	SELECT id
	FROM vendors 
	WHERE name = @vendor
)
INSERT INTO devices (snapshot_id, vendor_id, hostname, os_name, os_version, serial_number, management_ip)
SELECT s.id, v.id, @hostname, @os, @version, @serial, @ip
FROM snapshot_id s, vendor_id v;
`

	insertInterfaceQuery = `
WITH device_id AS (
	SELECT id
	FROM devices 
	ORDER BY id DESC
	LIMIT 1
)
INSERT INTO interfaces (device_id, name, mac, ip, mtu, bandwidth)
SELECT d.id, @name, @mac, @ip, @mtu, @bandwidth
FROM device_id d;
`
)

// SQL querie to get snapshot ids and timestamps.
const (
	selectTimestampsQuery = `
SELECT id, timestamp
FROM snapshots
ORDER BY timestamp DESC
LIMIT @limit;
`
)

// SQL querie to get a snapshot.
const (
	selectSnapshotQuery = `
SELECT
	s.timestamp,
	v.name AS vendor_name,
	d.id AS device_id,
	d.hostname,
	d.os_name,
	d.os_version,
	d.serial_number,
	d.management_ip,
	i.name AS interface_name,
	i.mac,
	i.ip,
	i.mtu,
	i.bandwidth
FROM
	snapshots s
	JOIN devices d ON s.id = d.snapshot_id
	JOIN vendors v ON v.id = d.vendor_id
	JOIN interfaces i ON d.id = i.device_id
WHERE
	s.id = @id
ORDER BY device_id ASC;
`
)

// SQL querie to delete a snapshot.
const (
	deleteSnapshotQuery = `
DELETE FROM snapshots
WHERE id = @id;
`
)

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
		createTableDevicesQuery,
		createTableInterfacesQuery,
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
