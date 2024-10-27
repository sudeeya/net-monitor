package postgresql

import (
	"net"
	"net/netip"
	"time"
)

type dbTimestamp struct {
	timestamp time.Time `db:"timestamp"`
}

type dbSnapshotPart struct {
	timestamp     time.Time        `db:"timestamp"`
	vendorName    string           `db:"vendor_name"`
	deviceID      int              `db:"device_id"`
	hostname      string           `db:"hostname"`
	osName        string           `db:"os_name"`
	osVersion     string           `db:"os_version"`
	serialNumber  string           `db:"serial_number"`
	managementIP  netip.Prefix     `db:"management_ip"`
	interfaceName string           `db:"interface_name"`
	mac           net.HardwareAddr `db:"mac"`
	ip            netip.Prefix     `db:"ip"`
	mtu           int64            `db:"mtu"`
	bandwidth     int64            `db:"bandwidth"`
}
