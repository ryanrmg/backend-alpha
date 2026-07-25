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

	ClearTradeIds(
		ctx context.Context,
	) error

	GetAllFillsOrdered(
		ctx context.Context,
	) ([]projectx.GatewayUserTrade, error)

	UpdateTradeId(
		ctx context.Context,
		fillId int,
		tradeId int64,
	) error
}
