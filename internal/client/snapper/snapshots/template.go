package snapshots

import "fmt"

const (
	nokiaVendor = "Nokia"
)

const (
	nokiaSRLinux = "nokia_srlinux"
)

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

type template struct {
	cmd     string
	file    string
	outputs []string
}

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

func getVendor(os string) (string, error) {
	switch os {
	case nokiaSRLinux:
		return nokiaVendor, nil
	default:
		return "", fmt.Errorf("unknown operating system: %s", os)
	}
}

func getTemplates(os string) ([]template, error) {
	switch os {
	case nokiaSRLinux:
		return nokiaSRLinuxTemplates, nil
	default:
		return nil, fmt.Errorf("unknown operating system: %s", os)
	}
}
