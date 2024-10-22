package converter

import (
	"net"

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
		Timestamp: model.Timestamp(snapshot.Timestamp.AsTime()),
		Devices:   devices,
	}, nil
}

func ToDeviceFromProto(device *pb.Snapshot_Device) (*model.Device, error) {
	ifaces := make([]model.Interface, len(device.Interfaces))

	for _, iface := range device.Interfaces {
		i, err := ToInterfaceFromProto(iface)
		if err != nil {
			return nil, err
		}

		ifaces = append(ifaces, *i)
	}

	return &model.Device{
		Vendor:     device.Vendor,
		Interfaces: ifaces,
	}, nil
}

func ToInterfaceFromProto(iface *pb.Snapshot_Device_Interface) (*model.Interface, error) {
	ip, ipnet, err := net.ParseCIDR(iface.IpAddress)
	if err != nil {
		return nil, err
	}

	return &model.Interface{
		Name:      iface.Name,
		IP:        ip,
		IPNet:     ipnet,
		MTU:       iface.Mtu,
		Bandwidth: iface.Bandwidth,
	}, nil
}
