package routes

import (
	"magnolia-test-backend/internal/handler"
	"net/http"
)

func RegisterFileRoutes(
	mux *http.ServeMux,
	handler *handler.FileHandler,
) {
	mux.HandleFunc(
		"POST /files",
		handler.Create,
	)

	mux.HandleFunc(
		"GET /files",
		handler.GetAll,
	)

	mux.HandleFunc(
		"GET /files/{id}",
		handler.GetByID,
	)

	mux.HandleFunc(
		"DELETE /files/{id}",
		handler.Delete,
	)
}
