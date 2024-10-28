package snapper

import (
	"github.com/sudeeya/net-monitor/internal/pkg/model"
)

type Snapper interface {
	Snap() *model.Snapshot
}
