package repository

import (
	"context"
	"time"

	projectx "github.com/ryanrmg/projectx-api"
)

type UserTradesJournalEntry struct {
	TradeId    int    `json:"tradeId"`
	AccountId  int    `json:"accountId"`
	ContractId string `json:"contractId"`

	EntryTimestamp string `json:"entryTimestamp"`
	ExitTimestamp  string `json:"exitTimestamp"`

	EntryPrice float64 `json:"entryPrice"`
	ExitPrice  float64 `json:"exitPrice"`

	EntrySize int `json:"entrySize"`
	ExitSize  int `json:"exitSize"`

	ProfitAndLoss float64 `json:"profitAndLoss"`
	Fees          float64 `json:"fees"`

	JournalNotes string `json:"journalNotes"`
}

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
