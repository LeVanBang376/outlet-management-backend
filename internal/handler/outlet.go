package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/response"
	"magnolia-test-backend/internal/service"
)

type OutletHandler struct {
	service *service.OutletService
}

func NewOutletHandler(service *service.OutletService) *OutletHandler {
	return &OutletHandler{
		service: service,
	}
}

func (h *OutletHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOutletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, "Created successfully", res)
}

func (h *OutletHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	outlet, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, customerrors.OutletErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, "Success", outlet)
}

func (h *OutletHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	pagination := response.NewPagination(r)

	outlets, err := h.service.GetAll(r.Context(), pagination)
	if err != nil {
		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.PaginatedJSON(w, http.StatusOK, "Success", outlets, pagination)
}

func (h *OutletHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var req dto.UpdateOutletRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Update(r.Context(), uint(id), req); err != nil {
		if errors.Is(err, customerrors.OutletErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NonDataJSON(w, http.StatusOK, "Updated successfully")
}

func (h *OutletHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		response.NonDataJSON(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := h.service.Delete(r.Context(), uint(id)); err != nil {
		if errors.Is(err, customerrors.OutletErrNotFound) {
			response.NonDataJSON(w, http.StatusNotFound, "Outlet not found")
			return
		}

		response.NonDataJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.NonDataJSON(w, http.StatusOK, "Deleted successfully")
}
