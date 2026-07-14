package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	projectx "github.com/ryanrmg/projectx-api"
)

type PostgresTradeRepository struct {
	store *DBStore
}

func NewPostgresTradeRepository(
	pool *pgxpool.Pool,
) *PostgresTradeRepository {

	return &PostgresTradeRepository{
		store: NewDBStore(pool),
	}
}

func (r *PostgresTradeRepository) GetTradesByAccount(
	ctx context.Context,
	accountId int,
) ([]projectx.GatewayUserTrade, error) {

	return r.store.GetTradesByAccount(ctx, accountId)
}

func (r *PostgresTradeRepository) SaveUserTrade(
	ctx context.Context,
	trade projectx.GatewayUserTrade,
) error {

	return r.store.SaveUserTrade(ctx, trade)
}

func (r *PostgresTradeRepository) CreateUserTradeTable(
	ctx context.Context,
) error {

	return r.store.CreateUserTradeTable(ctx)
}

func (r *PostgresTradeRepository) DeleteUserTable(
	ctx context.Context,
) error {

	return r.store.DeleteUserTable(ctx)
}

func (r *PostgresTradeRepository) GetLatestTradeTimestamp(
	ctx context.Context,
) (time.Time, error) {
	return r.store.GetLatestTradeTimestamp(ctx)
}
