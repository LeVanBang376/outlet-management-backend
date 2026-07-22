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

type FileHandler struct {
	service *service.FileService
}

func NewFileHandler(service *service.FileService) *FileHandler {
	return &FileHandler{
		service: service,
	}
}

func (h *FileHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.NonDataJSON(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}

	res, err := h.service.Create(r.Context(), req)
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
		http.StatusCreated,
		"Created successfully",
		res,
	)
}

func (h *FileHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	file, err := h.service.GetByID(
		r.Context(),
		uint(id),
	)
	if err != nil {
		if errors.Is(err, customerrors.FileErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"File not found",
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
		file,
	)
}

func (h *FileHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	files, err := h.service.GetAll(r.Context())
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
		files,
	)
}

func (h *FileHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		if errors.Is(err, customerrors.FileErrNotFound) {
			response.NonDataJSON(
				w,
				http.StatusNotFound,
				"File not found",
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
