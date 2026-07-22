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

type WorkingScheduleHandler struct {
	service *service.WorkingScheduleService
}

func NewWorkingScheduleHandler(service *service.WorkingScheduleService) *WorkingScheduleHandler {
	return &WorkingScheduleHandler{
		service: service,
	}
}

func (h *WorkingScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWorkingScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	if res, err := h.service.Create(r.Context(), req); err != nil {
		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	} else {
		response.JSON(
			w,
			http.StatusCreated,
			"Created successfully",
			res,
		)
		return
	}
}

func (h *WorkingScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}

	schedule, err := h.service.GetByID(
		r.Context(),
		uint(id),
	)

	if err != nil {
		if errors.Is(err, customerrors.WorkingScheduleErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"Working schedule not found",
			)
			return
		}

		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Success",
		schedule,
	)
}

func (h *WorkingScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.service.GetAll(r.Context())

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Success",
		schedules,
	)
}

func (h *WorkingScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}

	var req dto.UpdateWorkingScheduleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	res, err := h.service.Update(
		r.Context(),
		uint(id),
		req,
	)

	if err != nil {
		if errors.Is(err, customerrors.WorkingScheduleErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"Working schedule not found",
			)
			return
		}

		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		"Updated successfully",
		res,
	)
}

func (h *WorkingScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}

	err = h.service.Delete(
		r.Context(),
		uint(id),
	)

	if err != nil {
		if errors.Is(err, customerrors.WorkingScheduleErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"Working schedule not found",
			)
			return
		}

		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.NonDataJSON(
		w,
		http.StatusOK,
		"Deleted successfully",
	)
}
