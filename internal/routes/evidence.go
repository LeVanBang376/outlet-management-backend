package routes

import (
	"magnolia-test-backend/internal/handler"
	"net/http"
)

func RegisterEvidenceRoutes(
	mux *http.ServeMux,
	handler *handler.EvidenceHandler,
) {

	mux.HandleFunc(
		"POST /evidences",
		handler.Create,
	)

	mux.HandleFunc(
		"GET /evidences/{id}",
		handler.GetByID,
	)

	mux.HandleFunc(
		"GET /evidences",
		handler.GetByScheduleID,
	)

	mux.HandleFunc(
		"DELETE /evidences/{id}",
		handler.Delete,
	)
}
