// Package model defines basic data types.
package model

import (
	"net/netip"
	"time"
)

// Snapshot describes a snapshot.
type Snapshot struct {
	// Snapshot id.
	ID int `json:"id"`

	// The time at which the snapshot was created.
	Timestamp time.Time `json:"timestamp"`

	// A list of devices captured by the snapshot.
	Devices []Device `json:"devices"`
}

// Device describes a network device.
type Device struct {
	Hostname             string      `json:"hostname"`
	Vendor               string      `json:"vendor"`
	OSName               string      `json:"os_name"`
	OSVersion            string      `json:"os_version"`
	Serial               string      `json:"serial_number"`
	IsSnapshotSuccessful bool        `json:"is_snapshot_successful"`
	Interfaces           []Interface `json:"interfaces"`
}

// Interface describes a network device interface.
type Interface struct {
	Name string       `json:"name"`
	IsUp bool         `json:"is_up"`
	IP   netip.Prefix `json:"ip"`
	MTU  int64        `json:"mtu"`
}
