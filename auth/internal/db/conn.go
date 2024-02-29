package db

import (
	"context"

	"github.com/Foedie/foedie-server-v2/auth/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	DB *pgxpool.Pool
}

func NewDatabase(ctx context.Context, conf config.Config) (*Database, error) {
	sqlDriver, err := pgxpool.New(ctx, conf.DBUrl)
	if err != nil {
		return nil, err
	}
	return &Database{
		DB: sqlDriver,
	}, nil
}
