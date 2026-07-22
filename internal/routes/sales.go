package routes

import (
	"magnolia-test-backend/internal/handler"
	"net/http"
)

func RegisterSalesRoutes(
	mux *http.ServeMux,
	handler *handler.SalesHandler,
) {
	mux.HandleFunc(
		"GET /sales",
		handler.GetAll,
	)
}
