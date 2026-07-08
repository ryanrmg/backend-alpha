package repository

import (
    "context"

    projectx "github.com/ryanrmg/projectx-api"
)

type TradeRepository interface {
    GetTradesByAccount(
        ctx context.Context,
        accountID int,
    ) ([]projectx.GatewayUserTrade, error)
}