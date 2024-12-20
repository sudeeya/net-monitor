// Package converter provides functions to convert protobuf data to model data and vice versa.
package converter

import (
	"net"
	"net/netip"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
)

var errorStringPrefixNoSlash = "no '/'"

// ToProtoFromSnapshot converts model representation of snapshot to protobuf.
func ToProtoFromSnapshot(snapshot *model.Snapshot) *pb.Snapshot {
	devices := make([]*pb.Snapshot_Device, len(snapshot.Devices))

	for deviceIdx, device := range snapshot.Devices {
		devices[deviceIdx] = ToProtoFromDevice(device)
	}

	return &pb.Snapshot{
		Timestamp: timestamppb.New(snapshot.Timestamp),
		Devices:   devices,
	}
}

// ToProtoFromDevice converts model representation of device to protobuf.
func ToProtoFromDevice(device model.Device) *pb.Snapshot_Device {
	ifaces := make([]*pb.Snapshot_Device_Interface, len(device.Interfaces))

	for ifaceIdx, iface := range device.Interfaces {
		ifaces[ifaceIdx] = ToProtoFromInterface(iface)
	}

	return &pb.Snapshot_Device{
		Hostname:     device.Hostname,
		Vendor:       device.Vendor,
		OsName:       device.OSName,
		OsVersion:    device.OSVersion,
		Serial:       device.Serial,
		ManagementIp: netip.Prefix(device.ManagementIP).String(),
		Interfaces:   ifaces,
	}
}

// ToProtoFromInterface converts model representation of interface to protobuf.
func ToProtoFromInterface(iface model.Interface) *pb.Snapshot_Device_Interface {
	return &pb.Snapshot_Device_Interface{
		Name:      iface.Name,
		Mac:       net.HardwareAddr(iface.MAC).String(),
		Ip:        netip.Prefix(iface.IP).String(),
		Mtu:       iface.MTU,
		Bandwidth: iface.Bandwidth,
	}
}

// ToDeviceFromProto converts protobuf representation of snapshot to model.
func ToSnapshotFromProto(snapshot *pb.Snapshot) (*model.Snapshot, error) {
	devices := make([]model.Device, len(snapshot.Devices))

	for deviceIdx, device := range snapshot.Devices {
		d, err := ToDeviceFromProto(device)
		if err != nil {
			return nil, err
		}

		devices[deviceIdx] = *d
	}

	return &model.Snapshot{
		Timestamp: snapshot.Timestamp.AsTime(),
		Devices:   devices,
	}, nil
}

// ToDeviceFromProto converts protobuf representation of device to model.
func ToDeviceFromProto(device *pb.Snapshot_Device) (*model.Device, error) {
	ifaces := make([]model.Interface, len(device.Interfaces))

	managementIP, err := netip.ParsePrefix(device.ManagementIp)
	if err != nil {
		if err.Error() != errorStringPrefixNoSlash {
			return nil, err
		}
	}

	for ifaceIdx, iface := range device.Interfaces {
		i, err := ToInterfaceFromProto(iface)
		if err != nil {
			return nil, err
		}

		ifaces[ifaceIdx] = *i
	}

	return &model.Device{
		Hostname:     device.Hostname,
		Vendor:       device.Vendor,
		OSName:       device.OsName,
		OSVersion:    device.OsVersion,
		Serial:       device.Serial,
		ManagementIP: managementIP,
		Interfaces:   ifaces,
	}, nil
}

// ToInterfaceFromProto converts protobuf representation of interface to model.
func ToInterfaceFromProto(iface *pb.Snapshot_Device_Interface) (*model.Interface, error) {
	var mac net.HardwareAddr
	if iface.Mac != "" {
		m, err := net.ParseMAC(iface.Mac)
		if err != nil {
			return nil, err
		}
		mac = m
	}

	ip, err := netip.ParsePrefix(iface.Ip)
	if err != nil {
		if err.Error() != errorStringPrefixNoSlash {
			return nil, err
		}
	}

	return &model.Interface{
		Name:      iface.Name,
		MAC:       mac,
		IP:        ip,
		MTU:       iface.Mtu,
		Bandwidth: iface.Bandwidth,
	}, nil
}
