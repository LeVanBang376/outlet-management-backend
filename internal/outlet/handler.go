package outlet

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOutletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(r.Context(), req); err != nil {
		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NonDataJSON(w, http.StatusCreated, "Created successfully")
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	outlet, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "Success", outlet)
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	pagination := response.NewPagination(r)

	outlets, err := h.service.GetAll(r.Context(), pagination)
	if err != nil {
		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.PaginatedJSON(w, http.StatusOK, "Success", outlets, pagination)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req UpdateOutletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Update(r.Context(), uint(id), req); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NonDataJSON(w, http.StatusOK, "Updated successfully")
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, ErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NonDataJSON(w, http.StatusOK, "Deleted successfully")
}
