package service

import (
	"context"

	"github.com/ryanrmg/backend-alpha/internal/repository"
	projectx "github.com/ryanrmg/projectx-api"
)

type TradeService struct {
	repo repository.TradeRepository
}

func NewTradeService(
	repo repository.TradeRepository,
) *TradeService {

	return &TradeService{
		repo: repo,
	}
}

func (s *TradeService) GetTrades(
	ctx context.Context,
	accountID int,
) ([]projectx.GatewayUserTrade, error) {

	return s.repo.GetTradesByAccount(ctx, accountID)
}
