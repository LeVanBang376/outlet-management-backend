package handler

import (
	"net/http"

	"magnolia-test-backend/internal/response"
	"magnolia-test-backend/internal/service"
)

type SalesHandler struct {
	service *service.SalesService
}

func NewSalesHandler(service *service.SalesService) *SalesHandler {
	return &SalesHandler{
		service: service,
	}
}

func (h *SalesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
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
