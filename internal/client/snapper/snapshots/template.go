package snapshots

import "fmt"

const (
	ciscoVendor = "CISCO"
)

const (
	ciscoIOS = "cisco_ios"
)

const (
	versionOutput      = "VERSION"
	serialOutput       = "SERIAL"
	managementIPOutput = "MANAGEMENT_IP"
	interfaceOutput    = "INTERFACE"
	macAddressOutput   = "MAC_ADDRESS"
	ipAddressOutput    = "IP_ADDRESS"
	mtuOutput          = "MTU"
	bandwidthOutput    = "BANDWIDTH"
)

type template struct {
	cmd     string
	file    string
	outputs []string
}

var (
	ciscoIOSTemplates = []template{
		{
			cmd:  "show version",
			file: "templates/cisco_ios_show_version.textfsm",
			outputs: []string{
				versionOutput,
				serialOutput,
			},
		},
		{
			cmd:  "show interfaces",
			file: "templates/cisco_ios_show_interfaces.textfsm",
			outputs: []string{
				interfaceOutput,
				macAddressOutput,
				ipAddressOutput,
				mtuOutput,
				bandwidthOutput,
			},
		},
	}
)

func getVendor(os string) (string, error) {
	switch os {
	case ciscoIOS:
		return ciscoVendor, nil
	default:
		return "", fmt.Errorf("unknown operating system: %s", os)
	}
}

func getTemplates(os string) ([]template, error) {
	switch os {
	case ciscoIOS:
		return ciscoIOSTemplates, nil
	default:
		return nil, fmt.Errorf("unknown operating system: %s", os)
	}
}
