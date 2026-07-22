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

type EvidenceHandler struct {
	service *service.EvidenceService
}

func NewEvidenceHandler(service *service.EvidenceService) *EvidenceHandler {
	return &EvidenceHandler{
		service: service,
	}
}

func (h *EvidenceHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateEvidenceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	if err := h.service.Create(
		r.Context(),
		req,
	); err != nil {

		response.NonDataJSON(
			w,
			http.StatusInternalServerError,
			err.Error(),
		)

		return
	}

	response.NonDataJSON(
		w,
		http.StatusCreated,
		"Created successfully",
	)
}

func (h *EvidenceHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

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

	evidence, err := h.service.GetByID(
		r.Context(),
		uint(id),
	)

	if err != nil {
		if errors.Is(err, customerrors.EvidenceErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"Evidence not found",
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
		evidence,
	)
}

func (h *EvidenceHandler) GetByScheduleID(
	w http.ResponseWriter,
	r *http.Request,
) {

	scheduleID, err := strconv.ParseUint(
		r.URL.Query().Get("schedule_id"),
		10,
		64,
	)

	if err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid schedule_id",
		)

		return
	}

	evidences, err := h.service.GetByScheduleID(
		r.Context(),
		uint(scheduleID),
	)

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
		evidences,
	)
}

func (h *EvidenceHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

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

		if errors.Is(err, customerrors.EvidenceErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"Evidence not found",
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
