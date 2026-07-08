package repository

import (
    "context"

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
    accountID int,
) ([]projectx.GatewayUserTrade, error) {

    return r.store.GetTradesByAccount(ctx, accountID)
}