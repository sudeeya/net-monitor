package converter

import (
	"net"
	"net/netip"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
)

func ToSnapshotFromProto(snapshot *pb.Snapshot) (*model.Snapshot, error) {
	devices := make([]model.Device, len(snapshot.Devices))

	for _, device := range snapshot.Devices {
		d, err := ToDeviceFromProto(device)
		if err != nil {
			return nil, err
		}

		devices = append(devices, *d)
	}

	return &model.Snapshot{
		Timestamp: snapshot.Timestamp.AsTime(),
		Devices:   devices,
	}, nil
}

func ToDeviceFromProto(device *pb.Snapshot_Device) (*model.Device, error) {
	ifaces := make([]model.Interface, len(device.Interfaces))

	managementIP, err := netip.ParsePrefix(device.ManagementIp)
	if err != nil {
		return nil, err
	}

	for _, iface := range device.Interfaces {
		i, err := ToInterfaceFromProto(iface)
		if err != nil {
			return nil, err
		}

		ifaces = append(ifaces, *i)
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

func ToInterfaceFromProto(iface *pb.Snapshot_Device_Interface) (*model.Interface, error) {
	mac, err := net.ParseMAC(iface.Mac)
	if err != nil {
		return nil, err
	}

	ip, err := netip.ParsePrefix(iface.Ip)
	if err != nil {
		return nil, err
	}

	return &model.Interface{
		Name:      iface.Name,
		MAC:       mac,
		IP:        ip,
		MTU:       iface.Mtu,
		Bandwidth: iface.Bandwidth,
	}, nil
}
