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

func (r *PostgresTradeRepository) SaveUserFill(
	ctx context.Context,
	trade projectx.GatewayUserTrade,
	tradeId *int64,
) error {

	return r.store.SaveUserFill(ctx, trade, tradeId)
}

func (r *PostgresTradeRepository) CreateUserFillsTable(
	ctx context.Context,
) error {

	return r.store.CreateUserFillsTable(ctx)
}

func (r *PostgresTradeRepository) DeleteUserFillsTable(
	ctx context.Context,
) error {

	return r.store.DeleteUserFillsTable(ctx)
}

func (r *PostgresTradeRepository) GetLatestFill(
	ctx context.Context,
) (time.Time, *int64, error) {
	return r.store.GetLatestFill(ctx)
}
