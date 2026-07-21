package service

import (
	"sort"

	projectx "github.com/ryanrmg/projectx-api"
)

type Fill struct {
	projectx.GatewayUserTrade
	TradeID int64
}

type UserTrade struct {
	AccountID      int
	ContractID     string
	EntryTimestamp string
	ExitTimestamp  *string
}

func BuildTrades(
	fills []projectx.GatewayUserTrade,
) ([]Fill, []UserTrade) {

	sort.Slice(fills, func(i, j int) bool {
		return fills[i].CreationTimestamp < fills[j].CreationTimestamp
	})

	var result []Fill
	var trades []UserTrade

	var position int
	var currentTrade int64 = -1

	for _, fill := range fills {

		if position == 0 {

			currentTrade = int64(len(trades) + 1)

			trades = append(trades, UserTrade{
				AccountID:      fill.AccountId,
				ContractID:     fill.ContractId,
				EntryTimestamp: fill.CreationTimestamp,
			})
		}

		result = append(result, Fill{
			GatewayUserTrade: fill,
			TradeID:          currentTrade,
		})

		if fill.Side == 0 {
			position += fill.Size
		} else {
			position -= fill.Size
		}

		if position == 0 {

			exit := fill.CreationTimestamp

			trades[currentTrade-1].ExitTimestamp = &exit
		}
	}

	return result, trades
}
