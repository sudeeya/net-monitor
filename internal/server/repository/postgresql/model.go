package postgresql

import (
	"net"
	"net/netip"
	"time"
)

type dbTimestamp struct {
	ID        int       `db:"id"`
	Timestamp time.Time `db:"timestamp"`
}

type dbSnapshotPart struct {
	Timestamp     time.Time        `db:"timestamp"`
	VendorName    string           `db:"vendor_name"`
	DeviceID      int              `db:"device_id"`
	Hostname      string           `db:"hostname"`
	OSName        string           `db:"os_name"`
	OSVersion     string           `db:"os_version"`
	SerialNumber  string           `db:"serial_number"`
	ManagementIP  netip.Prefix     `db:"management_ip"`
	InterfaceName string           `db:"interface_name"`
	MAC           net.HardwareAddr `db:"mac"`
	IP            netip.Prefix     `db:"ip"`
	MTU           int64            `db:"mtu"`
	Bandwidth     int64            `db:"bandwidth"`
}
