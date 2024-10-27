package model

import (
	"net"
	"net/netip"
	"time"
)

type ID int

type Snapshot struct {
	Timestamp time.Time
	Devices   []Device
}

type Device struct {
	Hostname     string
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
