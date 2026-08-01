package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCandleRepository struct {
	store *DBStore
}

func NewPostgresCandleRepository(
	pool *pgxpool.Pool,
) *PostgresCandleRepository {

	return &PostgresCandleRepository{
		store: NewDBStore(pool),
	}
}

func (r *PostgresCandleRepository) GetCandles(
	ctx context.Context,
	contractId string,
	timeframe string,
	start time.Time,
	end time.Time,
) ([]Candle, error) {

	return r.store.GetCandles(ctx, contractId, timeframe, start, end)

}

func (r *PostgresCandleRepository) GetLatestCandle(
	ctx context.Context,
	contractId string,
) (time.Time, error) {

	return r.store.GetLatestCandle(ctx, contractId)

}

func (r *PostgresCandleRepository) SaveCandles(
	ctx context.Context,
	candles []Candle,
) error {

	return r.store.SaveCandles(ctx, candles)

}

