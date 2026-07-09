package api

import (
	"net/http"
)

func NewRouter(
	tradeHandler *TradeHandler,
) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc(
		"/api/user/trades",
		tradeHandler.GetTrades,
	)

	return mux
}
