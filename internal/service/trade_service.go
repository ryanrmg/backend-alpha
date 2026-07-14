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
) ([]projectx.GatewayUserTrade, error) {

	return s.repo.GetTradesByAccount(ctx, accountId)
}

func (s *TradeService) FetchTrades(
	ctx context.Context,
	accountId int,
) error {
	lastTradeTime, err := s.repo.GetLatestTradeTimestamp(ctx)
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

	for _, trade := range trades {
		if err := s.repo.SaveUserTrade(ctx, trade); err != nil {
			return err
		}
	}
	return nil
}
