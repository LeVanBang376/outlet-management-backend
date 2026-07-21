package evidence

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

func (h *Handler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateEvidenceRequest

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

func (h *Handler) GetByID(
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
		if errors.Is(err, ErrNotFound) {
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

func (h *Handler) GetByScheduleID(
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

func (h *Handler) Delete(
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

		if errors.Is(err, ErrNotFound) {
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
