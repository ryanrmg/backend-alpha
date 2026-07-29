package service

type MarketDataService struct {
	repo   CandleRepository
	client ProjectXClient
}

func (s *MarketDataService) GetCandles(
	ctx context.Context,
	contractId string,
	timeframe string,
	start time.Time,
	end time.Time,
) ([]Candle, error) {

}

func (s *MarketDataService) FetchCandles(
	ctx context.Context,
	contractId string,
) error {
	lastCandleTime, err := s.repo.GetLatestCandle(ctx)
	if err != nil {
		return err
	}

	startTime := lastCandleTime
	if startTime.IsZero() {
		// First sync for this account.
		startTime = time.Now().AddDate(0, 0, -30)
	}

	req := projectx.BarHistoryRequest{
		ContractId:        contractId,
		Live:              false,
		StartTime:         startTime.AddDate(0, 0, -2).UTC().Format(time.RFC3339), // check 2 days before it in case anything was missed
		EndTime:           time.Now().UTC().Format(time.RFC3339),
		Unit:              1, // always get the 1 minute for now, can derive
		UnitNumber:        1,
		Limit:             5000,
		IncludePartialBar: false,
	}

	candles, err := s.client.Market.History(ctx, req)
	if err != nil {
		return err
	}

	err = s.repo.SaveCandles(candles)
	if err != nil {
		return err
	}

	for _, fill := range fills {
		if err := s.repo.SaveUserFill(ctx, fill.Trade, fill.TradeId); err != nil {
			return err
		}
	}
	return nil
}
