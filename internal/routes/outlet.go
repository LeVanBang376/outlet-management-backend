package routes

import (
	"magnolia-test-backend/internal/handler"
	"net/http"
)

func RegisterOutletRoutes(
	mux *http.ServeMux,
	handler *handler.OutletHandler,
) {
	mux.HandleFunc("GET /outlets", handler.GetAll)
	mux.HandleFunc("GET /outlets/{id}", handler.GetByID)
	mux.HandleFunc("POST /outlets", handler.Create)
	mux.HandleFunc("PUT /outlets/{id}", handler.Update)
	mux.HandleFunc("DELETE /outlets/{id}", handler.Delete)
}
