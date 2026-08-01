package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ryanrmg/backend-alpha/internal/repository"
	projectx "github.com/ryanrmg/projectx-api"
)

type MarketDataService struct {
	repo   repository.CandleRepository
	client projectx.ProjectXClient
}

func (s *MarketDataService) GetCandles(
	ctx context.Context,
	contractId string,
	timeframe string,
	start time.Time,
	end time.Time,
) ([]repository.Candle, error) {
	return s.repo.GetCandles(ctx, contractId, timeframe, start, end)
}

func (s *MarketDataService) FetchCandles(
	ctx context.Context,
	contractId string,
) error {
	lastCandleTime, err := s.repo.GetLatestCandle(ctx, contractId)
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

	bars, err := s.client.Markets.History(ctx, req)
	if err != nil {
	    return err
	}

	candles := make([]repository.Candle, 0, len(bars))
	for _, bar := range bars {
	    ts, err := time.Parse(time.RFC3339, bar.Timestamp)
	    if err != nil {
	        return fmt.Errorf("invalid candle timestamp %q: %w", bar.Timestamp, err)
	    }

	    candles = append(candles, repository.Candle{
	        ContractId: contractId,
	        Timestamp:  ts,
	        Open:       bar.Open,
	        High:       bar.High,
	        Low:        bar.Low,
	        Close:      bar.Close,
	        Volume:     bar.Volume,
	    })
	}

	if err := s.repo.SaveCandles(ctx, candles); err != nil {
	    return err
	}


	return nil
}
