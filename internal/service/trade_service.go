package service

import (
	"context"
	"time"

	"github.com/ryanrmg/backend-alpha/internal/repository"
	projectx "github.com/ryanrmg/projectx-api"
)

type TradeService struct {
	repo   repository.TradeRepository
	client *projectx.ProjectXClient
}

type FillWithTradeId struct {
	Trade   projectx.GatewayUserTrade
	TradeId *int64
}

func NewTradeService(
	repo repository.TradeRepository,
	client *projectx.ProjectXClient,
) *TradeService {

	return &TradeService{
		repo:   repo,
		client: client,
	}
}

func (s *TradeService) GetTrades(
	ctx context.Context,
	accountId int,
) ([]repository.UserTradesJournalEntry, error) {

	return s.repo.GetTradesByAccount(ctx, accountId)
}

func (s *TradeService) FetchTrades(
	ctx context.Context,
	accountId int,
) error {
	lastTradeTime, lastTradeId, err := s.repo.GetLatestFill(ctx)
	if err != nil {
		return err
	}

	startTime := lastTradeTime
	if startTime.IsZero() {
		// First sync for this account.
		startTime = time.Now().AddDate(0, 0, -30)
	}

	req := projectx.TradeSearchRequest{
		AccountId:      accountId,
		StartTimestamp: startTime.UTC().Format(time.RFC3339),
		EndTimestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	trades, err := s.client.Trades.Search(ctx, req)
	if err != nil {
		return err
	}

	fills := s.AssignTrade(trades, lastTradeId)

	for _, fill := range fills {
		if err := s.repo.SaveUserFill(ctx, fill.Trade, fill.TradeId); err != nil {
			return err
		}
	}
	return nil
}

func (s *TradeService) AssignTrade(
	trades []projectx.GatewayUserTrade,
	lastTradeId *int64,
) []FillWithTradeId {

	fills := make([]FillWithTradeId, 0, len(trades))

	var (
		position       int
		currentTradeId int64
		nextTradeId    int64 = *lastTradeId
	)

	for _, trade := range trades {

		// Opening a new position starts a new trade.
		if position == 0 {
			currentTradeId = nextTradeId
			nextTradeId++
		}

		fills = append(fills, FillWithTradeId{
			Trade:   trade,
			TradeId: &currentTradeId,
		})

		// Update running position.
		if trade.Side == 0 {
			position += trade.Size
		} else {
			position -= trade.Size
		}

		// If position returns to zero, the next opening fill
		// will start a new trade.
		if position == 0 {
			currentTradeId = 0
		}
	}

	return fills
}
