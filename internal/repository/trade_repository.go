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

	GetLatestFillTimestamp(
		ctx context.Context,
	) (time.Time, error)

	SaveUserFill(
		ctx context.Context,
		trade projectx.GatewayUserTrade,
	) error
}
