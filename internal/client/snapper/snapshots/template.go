package snapshots

import "fmt"

// Vendors.
const (
	nokiaVendor = "Nokia"
)

// Operating systems.
const (
	nokiaSRLinux = "nokia_srlinux"
)

// Output data.
const (
	hostnameOutput     = "HOSTNAME"
	osOutput           = "OS"
	versionOutput      = "VERSION"
	serialOutput       = "SERIAL"
	managementIPOutput = "MANAGEMENT_IP"
	interfaceOutput    = "INTERFACE"
	macAddressOutput   = "MAC_ADDRESS"
	ipv4Output         = "IPV4"
	mtuOutput          = "MTU"
	bandwidthOutput    = "BANDWIDTH"
)

// template defines information needed to examine the configuration of a network device.
type template struct {
	// Command that need to be used on the device.
	cmd string

	// Name of the textfsm file needed to parse command response.
	file string

	// Output data present in the response.
	outputs []string
}

// OS-specific templates.
var (
	nokiaSRLinuxTemplates = []template{
		{
			cmd:  "show version",
			file: "templates/nokia_srlinux_show_version.textfsm",
			outputs: []string{
				hostnameOutput,
				osOutput,
				versionOutput,
			},
		},
		{
			cmd:  "show interface detail",
			file: "templates/nokia_srlinux_show_interface_detail.textfsm",
			outputs: []string{
				interfaceOutput,
				macAddressOutput,
				ipv4Output,
				mtuOutput,
				bandwidthOutput,
			},
		},
	}
)

// getVendor returns vendor by OS.
func getVendor(os string) (string, error) {
	switch os {
	case nokiaSRLinux:
		return nokiaVendor, nil
	default:
		return "", fmt.Errorf("unknown operating system: %s", os)
	}
}

// getTemplates returns templates by OS.
func getTemplates(os string) ([]template, error) {
	switch os {
	case nokiaSRLinux:
		return nokiaSRLinuxTemplates, nil
	default:
		return nil, fmt.Errorf("unknown operating system: %s", os)
	}
}
