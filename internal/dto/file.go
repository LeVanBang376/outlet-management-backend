package dto

import (
	"magnolia-test-backend/internal/model"
	"time"
)

type CreateFileRequest struct {
	ObjectKey   string `json:"object_key" validate:"required"`
	FileName    string `json:"file_name" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
	Size        int64  `json:"size" validate:"required,gt=0"`
}

type FileResponse struct {
	FileID      uint      `json:"file_id"`
	ObjectKey   string    `json:"object_key"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToFileResponse(file *model.File) *FileResponse {
	if file == nil {
		return nil
	}

	return &FileResponse{
		FileID:      file.FileID,
		ObjectKey:   file.ObjectKey,
		FileName:    file.FileName,
		ContentType: file.ContentType,
		Size:        file.Size,
		CreatedAt:   file.CreatedAt,
	}
}
