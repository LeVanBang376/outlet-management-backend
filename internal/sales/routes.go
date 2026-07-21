package sales

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	mux.HandleFunc(
		"GET /sales",
		handler.GetAll,
	)
}
