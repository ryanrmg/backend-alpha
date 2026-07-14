package repository

import (
	"context"
	"time"

	projectx "github.com/ryanrmg/projectx-api"
)

type TradeRepository interface {
	GetTradesByAccount(
		ctx context.Context,
		accountId int,
	) ([]projectx.GatewayUserTrade, error)

	GetLatestTradeTimestamp(
		ctx context.Context,
	) (time.Time, error)

	SaveUserTrade(
		ctx context.Context,
		trade projectx.GatewayUserTrade,
	) error
}
