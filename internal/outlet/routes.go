package outlet

import "net/http"

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
) {
	mux.HandleFunc("GET /outlets", handler.GetAll)
	mux.HandleFunc("GET /outlets/{id}", handler.GetByID)
	mux.HandleFunc("POST /outlets", handler.Create)
	mux.HandleFunc("PUT /outlets/{id}", handler.Update)
	mux.HandleFunc("DELETE /outlets/{id}", handler.Delete)
}
