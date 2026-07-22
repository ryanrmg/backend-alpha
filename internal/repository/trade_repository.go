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
	) ([]UserTradesJournalEntry, error)

	GetLatestFill(
		ctx context.Context,
	) (time.Time, *int64, error)

	SaveUserFill(
		ctx context.Context,
		trade projectx.GatewayUserTrade,
		tradeId *int64,
	) error
}
