package converter

import (
	"net"
	"net/netip"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProtoFromSnapshot(snapshot *model.Snapshot) *pb.Snapshot {
	devices := make([]*pb.Snapshot_Device, len(snapshot.Devices))

	for _, device := range snapshot.Devices {
		d := ToProtoFromDevice(device)
		devices = append(devices, d)
	}

	return &pb.Snapshot{
		Timestamp: timestamppb.New(snapshot.Timestamp),
		Devices:   devices,
	}
}

func ToProtoFromDevice(device model.Device) *pb.Snapshot_Device {
	ifaces := make([]*pb.Snapshot_Device_Interface, len(device.Interfaces))

	for _, iface := range device.Interfaces {
		i := ToProtoFromInterface(iface)
		ifaces = append(ifaces, i)
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

func ToProtoFromInterface(iface model.Interface) *pb.Snapshot_Device_Interface {
	return &pb.Snapshot_Device_Interface{
		Name:      iface.Name,
		Mac:       net.HardwareAddr(iface.MAC).String(),
		Ip:        netip.Prefix(iface.IP).String(),
		Mtu:       iface.MTU,
		Bandwidth: iface.Bandwidth,
	}
}

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
		ManagementIP: model.IPAddr(managementIP),
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
		MAC:       model.MACAddr(mac),
		IP:        model.IPAddr(ip),
		MTU:       iface.Mtu,
		Bandwidth: iface.Bandwidth,
	}, nil
}
