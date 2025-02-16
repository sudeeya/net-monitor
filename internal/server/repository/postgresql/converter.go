package postgresql

import (
	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

// toSnapshotFromDB creates a snapshot from a slice of database responses.
func toSnapshotFromDB(parts []dbSnapshotPart) model.Snapshot {
	if len(parts) == 0 {
		return model.Snapshot{}
	}

	deviceParts := make(map[int][]dbSnapshotPart, 1)
	for _, part := range parts {
		deviceParts[part.DeviceID] = append(deviceParts[part.DeviceID], part)
	}

	devices := make([]model.Device, len(deviceParts))
	for _, devicePart := range deviceParts {
		device := model.Device{
			Hostname:             devicePart[0].Hostname,
			Vendor:               devicePart[0].VendorName,
			OSName:               devicePart[0].OSName,
			OSVersion:            devicePart[0].OSVersion,
			Serial:               devicePart[0].SerialNumber,
			IsSnapshotSuccessful: devicePart[0].IsSnapshotSuccessful,
		}

		for _, part := range devicePart {
			iface := model.Interface{
				Name: part.InterfaceName,
				IsUp: part.IsUp,
				IP:   part.IP,
				MTU:  part.MTU,
			}
			device.Interfaces = append(device.Interfaces, iface)
		}
	}

	return model.Snapshot{
		Timestamp: parts[0].Timestamp,
		Devices:   devices,
	}
}
