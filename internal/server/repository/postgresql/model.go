package postgresql

import (
	"net/netip"
	"time"
)

// dbTimestamp is an auxiliary structure into which the database response is written.
type dbTimestamp struct {
	ID        int       `db:"id"`
	Timestamp time.Time `db:"timestamp"`
}

// dbSnapshotPart is an auxiliary structure into which the database response is written.
type dbSnapshotPart struct {
	Timestamp            time.Time    `db:"timestamp"`
	VendorName           string       `db:"vendor_name"`
	OSName               string       `db:"os_name"`
	OSVersion            string       `db:"os_version"`
	DeviceID             int          `db:"device_id"`
	Hostname             string       `db:"hostname"`
	SerialNumber         string       `db:"serial_number"`
	IsSnapshotSuccessful bool         `db:"is_snapshot_successful"`
	InterfaceName        string       `db:"interface_name"`
	IsUp                 bool         `db:"is_up"`
	IP                   netip.Prefix `db:"ip"`
	MTU                  int64        `db:"mtu"`
}
