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
		deviceParts[int(part.DeviceID.Int64)] = append(deviceParts[int(part.DeviceID.Int64)], part)
	}

	devices := make([]model.Device, len(deviceParts))
	devicesIdx := 0
	for _, devicePart := range deviceParts {
		device := model.Device{
			Hostname:             devicePart[0].Hostname.String,
			Vendor:               devicePart[0].VendorName.String,
			OSName:               devicePart[0].OSName.String,
			OSVersion:            devicePart[0].OSVersion.String,
			Serial:               devicePart[0].SerialNumber.String,
			IsSnapshotSuccessful: devicePart[0].IsSnapshotSuccessful.Bool,
		}

		for _, part := range devicePart {
			iface := model.Interface{
				Name: part.InterfaceName.String,
				IsUp: part.IsUp.Bool,
				IP:   part.IP,
				MTU:  part.MTU.Int64,
			}
			device.Interfaces = append(device.Interfaces, iface)
		}

		devices[devicesIdx] = device
		devicesIdx++
	}

	return model.Snapshot{
		Timestamp: parts[0].Timestamp.Time,
		Devices:   devices,
	}
}
