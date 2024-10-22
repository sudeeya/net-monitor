package model

import (
	"net"
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
	IP        net.IP
	IPNet     *net.IPNet
	MTU       int64
	Bandwidth int64
}
