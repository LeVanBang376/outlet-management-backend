package routes

import (
	"magnolia-test-backend/internal/handler"
	"net/http"
)

func RegisterWorkingScheduleRoutes(
	mux *http.ServeMux,
	handler *handler.WorkingScheduleHandler,
) {
	mux.HandleFunc(
		"POST /working-schedules",
		handler.Create,
	)

	mux.HandleFunc(
		"GET /working-schedules",
		handler.GetAll,
	)

	mux.HandleFunc(
		"GET /working-schedules/{id}",
		handler.GetByID,
	)

	mux.HandleFunc(
		"PUT /working-schedules/{id}",
		handler.Update,
	)

	mux.HandleFunc(
		"DELETE /working-schedules/{id}",
		handler.Delete,
	)
}
