package model

import (
	"net/netip"
	"time"
)

type ID int

type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Devices   []Device  `json:"devices"`
}

type Device struct {
	Hostname     string       `json:"hostname"`
	Vendor       string       `json:"vendor"`
	OSName       string       `json:"os_name"`
	OSVersion    string       `json:"os_version"`
	Serial       string       `json:"serial_number"`
	ManagementIP netip.Prefix `json:"management_ip"`
	Interfaces   []Interface  `json:"interfaces"`
}

type Interface struct {
	Name      string  `json:"name"`
	MAC       MACAddr `json:"mac"`
	IP        IPAddr  `json:"ip"`
	MTU       int64   `json:"mtu"`
	Bandwidth int64   `json:"bandwidth"`
}
