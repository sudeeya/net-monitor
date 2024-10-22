package model

import (
	"net/netip"
	"time"
)

type Timestamp time.Time

type Snapshot struct {
	Timestamp Timestamp
	Devices   []Device
}

type Device struct {
	Vendor     string
	Interfaces []Interface
}

type Interface struct {
	Name      string
	IPAddress netip.Prefix
	MTU       int64
	Bandwidth int64
}
