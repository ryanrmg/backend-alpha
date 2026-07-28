package api

func (h *MarketDataHandler) GetCandles(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	contractID := r.URL.Query().Get("contractId")
	timeframe := r.URL.Query().Get("timeframe")

	start, err := time.Parse(
		time.RFC3339,
		r.URL.Query().Get("start"),
	)
	if err != nil {
		http.Error(w, "invalid start", http.StatusBadRequest)
		return
	}

	end, err := time.Parse(
		time.RFC3339,
		r.URL.Query().Get("end"),
	)
	if err != nil {
		http.Error(w, "invalid end", http.StatusBadRequest)
		return
	}

	candles, err := h.service.GetCandles(
		ctx,
		contractID,
		timeframe,
		start,
		end,
	)
	if err != nil {
		http.Error(w, "failed to get candles", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(candles)
}
