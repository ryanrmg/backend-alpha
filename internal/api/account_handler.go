package api

import (
	"encoding/json"
	"github.com/ryanrmg/backend-alpha/internal/service"
	"log"
	"net/http"
)

type AccountHandler struct {
	service *service.AccountService
}

func NewAccountHandler(
	service *service.AccountService,
) *AccountHandler {

	return &AccountHandler{
		service: service,
	}
}

func (h *AccountHandler) GetAccounts(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Printf("Received %s %s from %s",
		r.Method,
		r.URL.String(),
		r.RemoteAddr,
	)

	accounts, err := h.service.GetAccounts(
		r.Context(),
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(accounts)
}
