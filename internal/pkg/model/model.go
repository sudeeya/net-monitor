package model

import "time"

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
	IPAddress string
	MTU       int64
	Bandwidth int64
}
