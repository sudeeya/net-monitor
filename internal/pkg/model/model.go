package model

import (
	"net"
	"net/netip"
	"time"
)

type Timestamp time.Time

type Snapshot struct {
	Timestamp Timestamp
	Devices   []Device
}

type Device struct {
	Vendor       string
	OSName       string
	OSVersion    string
	Serial       string
	ManagementIP netip.Prefix
	Interfaces   []Interface
}

type Interface struct {
	Name      string
	MAC       net.HardwareAddr
	IP        netip.Prefix
	MTU       int64
	Bandwidth int64
}
