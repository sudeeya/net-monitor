package postgresql

import (
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
)

// dbTimestamp is an auxiliary structure into which the database response is written.
type dbTimestamp struct {
	ID        pgtype.Int8        `db:"id"`
	Timestamp pgtype.Timestamptz `db:"timestamp"`
}

// dbSnapshotPart is an auxiliary structure into which the database response is written.
type dbSnapshotPart struct {
	ID                   pgtype.Int8        `db:"id"`
	Timestamp            pgtype.Timestamptz `db:"timestamp"`
	VendorName           pgtype.Text        `db:"vendor_name"`
	OSName               pgtype.Text        `db:"os_name"`
	OSVersion            pgtype.Text        `db:"os_version"`
	DeviceID             pgtype.Int8        `db:"device_id"`
	Hostname             pgtype.Text        `db:"hostname"`
	SerialNumber         pgtype.Text        `db:"serial_number"`
	IsSnapshotSuccessful pgtype.Bool        `db:"is_snapshot_successful"`
	InterfaceName        pgtype.Text        `db:"interface_name"`
	IsUp                 pgtype.Bool        `db:"is_up"`
	IP                   netip.Prefix       `db:"ip"`
	MTU                  pgtype.Int8        `db:"mtu"`
}
