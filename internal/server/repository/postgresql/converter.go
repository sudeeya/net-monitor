package postgresql

import (
	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

// toSnapshotFromDB creates a snapshot from a slice of database responses.
func toSnapshotFromDB(parts []dbSnapshotPart) model.Snapshot {
	if len(parts) == 0 {
		return model.Snapshot{}
	}

	devices := make([]model.Device, 0)
	currentID := 0
	for _, part := range parts {
		switch {
		case part.DeviceID != currentID:
			currentID = part.DeviceID
			ifaces := []model.Interface{
				{
					Name:      part.InterfaceName,
					MAC:       part.MAC,
					IP:        part.IP,
					MTU:       part.MTU,
					Bandwidth: part.Bandwidth,
				},
			}
			devices = append(devices, model.Device{
				Hostname:     part.Hostname,
				Vendor:       part.VendorName,
				OSName:       part.OSName,
				OSVersion:    part.OSVersion,
				Serial:       part.SerialNumber,
				ManagementIP: part.ManagementIP,
				Interfaces:   ifaces,
			})
		default:
			iface := model.Interface{
				Name:      part.InterfaceName,
				MAC:       part.MAC,
				IP:        part.IP,
				MTU:       part.MTU,
				Bandwidth: part.Bandwidth,
			}
			devices[len(devices)-1].Interfaces = append(devices[len(devices)-1].Interfaces, iface)
		}
	}

	return model.Snapshot{
		Timestamp: parts[0].Timestamp,
		Devices:   devices,
	}
}
