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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const limitInSeconds = 100

// Client describes client.
type Client struct {
	logger  *zap.Logger
	snapper snapper.Snapper
	conn    *grpc.ClientConn
	client  pb.SnapshotsClient
}

// NewClient returns client object.
func NewClient(
	logger *zap.Logger,
	snapper snapper.Snapper,
	serverAddr string,
) (*Client, error) {
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewSnapshotsClient(conn)

	return &Client{
		logger:  logger,
		snapper: snapper,
		conn:    conn,
		client:  client,
	}, nil
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

// Close tears down connections.
func (c *Client) Close() error {
	return c.conn.Close()
}
