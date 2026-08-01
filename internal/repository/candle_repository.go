package repository

import (
	"context"
	"time"
)

type Candle struct {
	ContractId string
	Timeframe  string

	Timestamp time.Time

	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

type CandleRepository interface {
	GetCandles(
		ctx context.Context,
		contractId string,
		timeframe string,
		start time.Time,
		end time.Time,
	) ([]Candle, error)

	SaveCandles(
		ctx context.Context,
		candles []Candle,
	) error

	GetLatestCandle(
		ctx context.Context,
		contractId string,
	) (time.Time, error)
}
