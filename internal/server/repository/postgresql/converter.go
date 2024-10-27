package postgresql

import (
	"time"

	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

func toDBFromTimestamp(t model.Timestamp) dbTimestamp {
	return dbTimestamp{
		timestamp: time.Time(t),
	}
}

func toTimestampFromDB(dbt dbTimestamp) model.Timestamp {
	return model.Timestamp(dbt.timestamp)
}

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
					MAC:       part.mac,
					IP:        part.ip,
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
				MAC:       part.mac,
				IP:        part.ip,
				MTU:       part.mtu,
				Bandwidth: part.bandwidth,
			}
			devices[len(devices)-1].Interfaces = append(devices[len(devices)-1].Interfaces, iface)
		}
	}

	return model.Snapshot{
		Timestamp: model.Timestamp(parts[0].timestamp),
		Devices:   devices,
	}
}
