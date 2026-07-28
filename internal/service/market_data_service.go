package service

type MarketDataService struct {
	repo   CandleRepository
	client ProjectXClient
}

func (s *MarketDataService) GetCandles(
	ctx context.Context,
	contractID string,
	timeframe string,
	start time.Time,
	end time.Time,
) ([]Candle, error) {

}
