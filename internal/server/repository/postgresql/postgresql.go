package postgresql

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sudeeya/net-monitor/internal/pkg/model"
	"github.com/sudeeya/net-monitor/internal/server/repository"
	"go.uber.org/zap"
)

const driverName string = "pgx"

var _ repository.Repository = (*postgreSQL)(nil)

type postgreSQL struct {
	logger *zap.Logger
	db     *sqlx.DB
}

func NewPostgreSQL(logger *zap.Logger, dsn string) (*postgreSQL, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	return &postgreSQL{
		logger: logger,
		db:     db,
	}, nil
}

func (p *postgreSQL) DeleteSnapshot(ctx context.Context, timestamp model.Timestamp) error {
	panic("unimplemented")
}

func (p *postgreSQL) GetSnapshot(ctx context.Context, timestamp model.Timestamp) (model.Snapshot, error) {
	panic("unimplemented")
}

func (p *postgreSQL) StoreSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	panic("unimplemented")
}
