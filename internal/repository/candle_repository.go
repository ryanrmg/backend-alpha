package repository

type Candle struct {
	ContractID string
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
		contractID string,
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
		contractID string,
		timeframe string,
	) (*Candle, error)
}
