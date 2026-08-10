package http

import (
	"net/http"

	"github.com/Aifyel/petCryptoCurrency/internal/reader/usecase"
)

type Handler struct {
	Handler usecase.RateQuery
}

func newHandler(handler usecase.RateQuery) *Handler {
	return &Handler{Handler: handler}
}

func (h *Handler) GetStatistics(w http.ResponseWriter, r *http.Request) {
}
