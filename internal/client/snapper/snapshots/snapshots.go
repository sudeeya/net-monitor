package snapshots

import (
	"encoding/json"
	"errors"
	"net"
	"net/netip"
	"os"
	"strconv"
	"time"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/sudeeya/net-monitor/internal/client/snapper"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"go.uber.org/zap"
)

var _ snapper.Snapper = (*snapshots)(nil)

type snapshots struct {
	logger  *zap.Logger
	targets []target
}

type target struct {
	hostname  string
	vendor    string
	os        string
	driver    *generic.Driver
	templates []template
}

type targetConfig struct {
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSnapshots(logger *zap.Logger, targetsFile string) (*snapshots, error) {
	cfgs, err := extractConfigs(targetsFile)
	if err != nil {
		return nil, err
	}

	targets, err := formTargets(cfgs)
	if err != nil {
		return nil, err
	}

	return &snapshots{
		logger:  logger,
		targets: targets,
	}, nil
}

func (s *snapshots) Snap() (*model.Snapshot, error) {
	timestamp := time.Now()
	devices := make([]model.Device, 0)

	connErrs := error(nil)

	for _, t := range s.targets {
		err := t.driver.Open()
		if err != nil {
			connErrs = errors.Join(connErrs, err)
			continue
		}
		defer t.driver.Close()

		device := model.Device{
			Hostname: t.hostname,
			Vendor:   t.vendor,
			OSName:   t.os,
		}

		ifaces := make([]model.Interface, 0)

		for _, template := range t.templates {
			response, err := t.driver.SendCommand(template.cmd)
			if err != nil {
				return nil, err
			}

			parsed, err := response.TextFsmParse(template.file)
			if err != nil {
				return nil, err
			}

			for _, p := range parsed {
				iface := model.Interface{}

				for _, output := range template.outputs {
					if p == nil {
						continue
					}

					value, ok := p[output].(string)
					if !ok {
						continue
					}
					if value == "" {
						continue
					}

					switch output {
					case versionOutput:
						device.OSVersion = value
					case serialOutput:
						device.Serial = value
					case managementIPOutput:
						ip, err := netip.ParsePrefix(value)
						if err != nil {
							return nil, err
						}
						device.ManagementIP = ip
					case interfaceOutput:
						iface.Name = value
					case macAddressOutput:
						mac, err := net.ParseMAC(value)
						if err != nil {
							return nil, err
						}
						iface.MAC = mac
					case ipv4Output:
						ip, err := netip.ParsePrefix(value)
						if err != nil {
							return nil, err
						}
						iface.IP = ip
					case mtuOutput:
						mtu, err := strconv.Atoi(value)
						if err != nil {
							return nil, err
						}
						iface.MTU = int64(mtu)
					case bandwidthOutput:
						bandwidth, err := strconv.Atoi(value)
						if err != nil {
							return nil, err
						}
						iface.Bandwidth = int64(bandwidth)
					}
				}

				if iface.Name != "" {
					ifaces = append(ifaces, iface)
				}
			}
		}

		device.Interfaces = ifaces

		devices = append(devices, device)
	}

	return &model.Snapshot{
		Timestamp: timestamp,
		Devices:   devices,
	}, connErrs
}

func extractConfigs(targetFile string) ([]targetConfig, error) {
	file, err := os.Open(targetFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfgs := make([]targetConfig, 0)

	if err := json.NewDecoder(file).Decode(&cfgs); err != nil {
		return nil, err
	}

	return cfgs, nil
}

func formTargets(cfgs []targetConfig) ([]target, error) {
	targets := make([]target, len(cfgs))

	for cfgIdx, cfg := range cfgs {
		driver, err := generic.NewDriver(
			cfg.Hostname,
			options.WithAuthNoStrictKey(),
			options.WithAuthUsername(cfg.Username),
			options.WithAuthPassword(cfg.Password),
		)
		if err != nil {
			return nil, err
		}

		vendor, err := getVendor(cfg.OS)
		if err != nil {
			return nil, err
		}

		templates, err := getTemplates(cfg.OS)
		if err != nil {
			return nil, err
		}

		targets[cfgIdx] = target{
			hostname:  cfg.Hostname,
			vendor:    vendor,
			os:        cfg.OS,
			driver:    driver,
			templates: templates,
		}
	}

	return targets, nil
}
