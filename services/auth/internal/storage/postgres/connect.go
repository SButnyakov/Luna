package postgres

import (
	"context"
	"fmt"

	"github.com/SButnyakov/luna/id/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg config.PG) (*pgxpool.Pool, error) {
	const op = "storage.postgres.connect.Connect"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: unable to connect to postgres database: %w", op, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: couldn't ping postgres database: %w", op, err)
	}

	return pool, nil
}
