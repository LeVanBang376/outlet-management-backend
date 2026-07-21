package sales

import (
	"net/http"

	"magnolia-test-backend/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	sales, err := h.service.GetAll(r.Context())

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Success",
		sales,
	)
}
