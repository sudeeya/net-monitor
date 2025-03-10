// Package snapshots defines object that creates snapshots by connecting to network devices via SSH.
package snapshots

import (
	"encoding/json"
	"net/netip"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/scrapli/scrapligo/driver/options"

	"github.com/sudeeya/net-monitor/internal/client/snapper"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

var _ snapper.Snapper = (*snapshots)(nil)

// snapshots implements the [Snapper] interface.
type snapshots struct {
	logger  *zap.Logger
	targets []target
}

// target defines a target device.
type target struct {
	cfg       targetConfig
	templates []template
}

// targetConfig defines device OS and information needed for an SSH connection.
type targetConfig struct {
	OS             string `json:"os"`
	Hostname       string `json:"hostname"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	PrivateKeyPath string `json:"private_key_path"`
	Passphrase     string `json:"passphrase"`
	NoStrictKey    bool   `json:"no_strict_key"`
}

// NewSnapshots returns snapshots object.
// The function extracts target network devices from  a json file.
func NewSnapshots(logger *zap.Logger, targetsFile string) (*snapshots, error) {
	logger.Sugar().Infof("Extracting configs from file %s", targetsFile)
	cfgs, err := extractConfigs(targetsFile)
	if err != nil {
		return nil, err
	}

	logger.Info("Forming a list of target devices")
	targets, err := formTargets(cfgs)
	if err != nil {
		return nil, err
	}

	return &snapshots{
		logger:  logger,
		targets: targets,
	}, nil
}

type snapResult struct {
	hostname string
	device   *model.Device
	err      error
}

// Snap implements the [Snapper] interface.
func (s *snapshots) Snap() (*model.Snapshot, error) {
	timestamp := time.Now()
	devices := make([]model.Device, 0)
	resultChan := make(chan snapResult, len(s.targets))

	for _, t := range s.targets {
		go func(t target) {
			device, err := s.snapTarget(t)
			resultChan <- snapResult{t.cfg.Hostname, device, err}
		}(t)
	}

	for range len(s.targets) {
		res := <-resultChan

		if res.err != nil {
			s.logger.Sugar().Errorf("Failed to snap target %s: %s", res.hostname, res.err.Error())
		} else {
			devices = append(devices, *res.device)
		}
	}

	return &model.Snapshot{
		Timestamp: timestamp,
		Devices:   devices,
	}, nil
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
		templates, err := getTemplates(cfg.OS)
		if err != nil {
			return nil, err
		}

		targets[cfgIdx] = target{
			cfg:       cfg,
			templates: templates,
		}
	}

	return targets, nil
}

func (s *snapshots) snapTarget(t target) (*model.Device, error) {
	vendor, err := getVendor(t.cfg.OS)
	if err != nil {
		return nil, err
	}

	driver, err := newTargetDriver(t)
	if err != nil {
		return nil, err
	}

	device := &model.Device{
		Hostname: t.cfg.Hostname,
		Vendor:   vendor,
		OSName:   t.cfg.OS,
	}

	s.logger.Sugar().Infof("Trying to connect to %s", t.cfg.Hostname)
	err = driver.Open()
	defer driver.Close()
	if err != nil {
		s.logger.Sugar().Errorf("Failed to connect to %s: %s", t.cfg.Hostname, err.Error())

		device.IsSnapshotSuccessful = false

		return device, nil
	}

	s.logger.Sugar().Infof("Connection to %s established", t.cfg.Hostname)

	ifaces := make([]model.Interface, 0)

	for _, template := range t.templates {
		s.logger.Sugar().Infof("Sending command: %s", template.cmd)
		response, err := driver.SendCommand(template.cmd)
		if err != nil {
			return nil, err
		}

		s.logger.Info("Parsing response")
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
				case hostnameOutput:
					device.Hostname = value
				case osOutput:
					device.OSName = value
				case versionOutput:
					device.OSVersion = value
				case serialOutput:
					device.Serial = value
				case interfaceOutput:
					iface.Name = value
				case stateOutput:
					switch value {
					case "up":
						iface.IsUp = true
					case "down":
						iface.IsUp = false
					}
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
				}
			}

			if iface.Name != "" {
				ifaces = append(ifaces, iface)
			}
		}
	}

	device.Interfaces = ifaces
	device.IsSnapshotSuccessful = true

	return device, nil
}

func newTargetDriver(t target) (*generic.Driver, error) {
	var (
		driver *generic.Driver
		err    error
	)

	switch {
	case t.cfg.PrivateKeyPath != "" && t.cfg.NoStrictKey:
		driver, err = generic.NewDriver(
			t.cfg.Hostname,
			options.WithAuthNoStrictKey(),
			options.WithAuthPrivateKey(t.cfg.PrivateKeyPath, t.cfg.Passphrase),
		)
		if err != nil {
			return nil, err
		}
	case t.cfg.PrivateKeyPath != "":
		driver, err = generic.NewDriver(
			t.cfg.Hostname,
			options.WithAuthPrivateKey(t.cfg.PrivateKeyPath, t.cfg.Passphrase),
		)
		if err != nil {
			return nil, err
		}
	case t.cfg.NoStrictKey:
		driver, err = generic.NewDriver(
			t.cfg.Hostname,
			options.WithAuthNoStrictKey(),
			options.WithAuthUsername(t.cfg.Username),
			options.WithAuthPassword(t.cfg.Password),
		)
		if err != nil {
			return nil, err
		}
	default:
		driver, err = generic.NewDriver(
			t.cfg.Hostname,
			options.WithAuthUsername(t.cfg.Username),
			options.WithAuthPassword(t.cfg.Password),
		)
		if err != nil {
			return nil, err
		}
	}

	return driver, nil
}
