package evidence

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
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
