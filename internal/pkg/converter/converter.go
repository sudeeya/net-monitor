package converter

import (
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
)

func ToSnapshotFromProto(snapshot *pb.Snapshot) *model.Snapshot {
	devices := make([]model.Device, len(snapshot.Devices))

	for _, device := range snapshot.Devices {
		devices = append(devices, *ToDeviceFromProto(device))
	}

	return &model.Snapshot{
		Timestamp: model.Timestamp(snapshot.Timestamp.AsTime()),
		Devices:   devices,
	}
}

func ToDeviceFromProto(device *pb.Snapshot_Device) *model.Device {
	ifaces := make([]model.Interface, len(device.Interfaces))

	for _, iface := range device.Interfaces {
		ifaces = append(ifaces, *ToInterfaceFromProto(iface))
	}

	return &model.Device{
		Vendor:     device.Vendor,
		Interfaces: ifaces,
	}
}

func ToInterfaceFromProto(iface *pb.Snapshot_Device_Interface) *model.Interface {
	return &model.Interface{
		Name:      iface.Name,
		MTU:       iface.Mtu,
		Bandwidth: iface.Bandwidth,
	}
}
