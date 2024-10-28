package snapshots

import (
	"github.com/scrapli/scrapligo/driver/generic"
	"github.com/sudeeya/net-monitor/internal/client/snapper"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"go.uber.org/zap"
)

var _ snapper.Snapper = (*snapshots)(nil)

type snapshots struct {
	logger  *zap.Logger
	targets []*generic.Driver
}

func NewSnapshots(logger *zap.Logger, targetsFile string) (*snapshots, error) {
	targets := make([]*generic.Driver, 0)

	return &snapshots{
		logger:  logger,
		targets: targets,
	}, nil
}

func (s *snapshots) Snap() *model.Snapshot {
	panic("unimplemented")
}
