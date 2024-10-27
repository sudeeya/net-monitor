package postgresql

import (
	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

func toSnapshotFromDB(parts []dbSnapshotPart) model.Snapshot {
	if len(parts) == 0 {
		return model.Snapshot{}
	}

	devices := make([]model.Device, 0)
	currentID := 0
	for _, part := range parts {
		switch {
		case part.deviceID != currentID:
			currentID = part.deviceID
			ifaces := []model.Interface{
				{
					Name:      part.interfaceName,
					MAC:       model.MACAddr(part.mac),
					IP:        model.IPAddr(part.ip),
					MTU:       part.mtu,
					Bandwidth: part.bandwidth,
				},
			}
			devices = append(devices, model.Device{
				Hostname:     part.hostname,
				Vendor:       part.vendorName,
				OSName:       part.osName,
				OSVersion:    part.osVersion,
				Serial:       part.serialNumber,
				ManagementIP: part.managementIP,
				Interfaces:   ifaces,
			})
		default:
			iface := model.Interface{
				Name:      part.interfaceName,
				MAC:       model.MACAddr(part.mac),
				IP:        model.IPAddr(part.ip),
				MTU:       part.mtu,
				Bandwidth: part.bandwidth,
			}
			devices[len(devices)-1].Interfaces = append(devices[len(devices)-1].Interfaces, iface)
		}
	}

	return model.Snapshot{
		Timestamp: parts[0].timestamp,
		Devices:   devices,
	}
}
