package api

import (
	"encoding/json"
	"github.com/ryanrmg/backend-alpha/internal/service"
	"log"
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

	log.Printf("Received %s %s from %s",
		r.Method,
		r.URL.String(),
		r.RemoteAddr,
	)
	accountId, err := strconv.Atoi(
		r.URL.Query().Get("accountId"),
	)

	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to fetch trades from remote")
	err = h.service.FetchTrades(r.Context(), accountId)
	if err != nil {
		log.Printf("Failed to get trades from remote", err)
	}

	trades, err := h.service.GetTrades(
		r.Context(),
		accountId,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(trades)
}
