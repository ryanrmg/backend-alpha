package api

import (
	"net/http"
)

func NewRouter(
	tradeHandler *TradeHandler,
	accountHandler *AccountHandler,
) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/user/trades",
		tradeHandler.GetTrades,
	)

	mux.HandleFunc(
		"api/user/accounts",
		accountHandler.GetAccounts,
	)

	return mux
}
