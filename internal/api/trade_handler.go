package api

import (
	"encoding/json"
	"github.com/ryanrmg/backend-alpha/internal/service"
	"net/http"
	"strconv"
)

type TradeHandler struct {
	service *service.TradeService
}

func NewTradeHandler(
	service *service.TradeService,
) *TradeHandler {

	return &TradeHandler{
		service: service,
	}
}

func (h *TradeHandler) GetTrades(
	w http.ResponseWriter,
	r *http.Request,
) {

	accountID, err := strconv.Atoi(
		r.URL.Query().Get("accountId"),
	)

	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	trades, err := h.service.GetTrades(
		r.Context(),
		accountID,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(trades)
}
