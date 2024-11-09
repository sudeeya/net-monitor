// Package model defines basic data types.
package model

import (
	"net"
	"net/netip"
	"time"
)

// ID describes object identifier.
type ID int

// Snapshot describes a snapshot.
type Snapshot struct {
	// The time at which the snapshot was created.
	Timestamp time.Time `json:"timestamp"`

	// A list of devices captured by the snapshot.
	Devices []Device `json:"devices"`
}

// Device describes a network device.
type Device struct {
	Hostname     string       `json:"hostname"`
	Vendor       string       `json:"vendor"`
	OSName       string       `json:"os_name"`
	OSVersion    string       `json:"os_version"`
	Serial       string       `json:"serial_number"`
	ManagementIP netip.Prefix `json:"management_ip"`
	Interfaces   []Interface  `json:"interfaces"`
}

// Interface describes a network device interface.
type Interface struct {
	Name      string           `json:"name"`
	MAC       net.HardwareAddr `json:"mac"`
	IP        netip.Prefix     `json:"ip"`
	MTU       int64            `json:"mtu"`
	Bandwidth int64            `json:"bandwidth"`
}
