// Package client defines object that interacts with [Snapper] and communicates with the server.
package client

import (
	"context"
	"fmt"
	"time"

	"github.com/sudeeya/net-monitor/internal/client/snapper"
	"github.com/sudeeya/net-monitor/internal/pkg/converter"
	"github.com/sudeeya/net-monitor/internal/pkg/pb"
	"go.uber.org/zap"
)

const limitInSeconds = 100

// Client describes client.
type Client struct {
	logger  *zap.Logger
	snapper snapper.Snapper
	client  pb.SnapshotsClient
}

// NewClient returns client object.
func NewClient(
	logger *zap.Logger,
	snapper snapper.Snapper,
	client pb.SnapshotsClient,
) *Client {
	return &Client{
		logger:  logger,
		snapper: snapper,
		client:  client,
	}
}

// UploadSnapshot requests the [Snapper] to make snapshot and sends it to the server.
func (c *Client) UploadSnapshot() error {
	ctx, cancel := context.WithTimeout(context.Background(), limitInSeconds*time.Second)
	defer cancel()

	c.logger.Info("Snapshot is being created")
	s, err := c.snapper.Snap()
	if err != nil {
		return err
	}
	c.logger.Info("Snapshot is ready to be saved")

	snapshot := converter.ToProtoFromSnapshot(s)

	response, err := c.client.SaveSnapshot(ctx, &pb.SaveSnapshotRequest{Snapshot: snapshot})
	if err != nil {
		return err
	}
	if response.Error != "" {
		return fmt.Errorf(response.Error)
	}
	c.logger.Info("Snapshot has been saved")

	return nil
}
