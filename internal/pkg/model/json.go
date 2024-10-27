package model

import (
	"encoding/json"
	"net"
	"net/netip"
)

type MACAddr net.HardwareAddr

type IPAddr netip.Prefix

func (mac MACAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(net.HardwareAddr(mac).String())
}

func (ip IPAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(netip.Prefix(ip).String())
}
