package dto

import (
	"magnolia-test-backend/internal/model"
	"time"
)

type CreateEvidenceRequest struct {
	ScheduleID  uint   `json:"schedule_id"`
	ObjectKey   string `json:"object_key"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type EvidenceResponse struct {
	EvidenceID  uint      `json:"evidence_id"`
	ScheduleID  uint      `json:"schedule_id"`
	ObjectKey   string    `json:"object_key"`
	FileName    string    `json:"file_name"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

func ToEvidenceResponse(evidence *model.Evidence) *EvidenceResponse {
	if evidence == nil {
		return nil
	}

	return &EvidenceResponse{
		EvidenceID:  evidence.EvidenceID,
		ScheduleID:  evidence.ScheduleID,
		ObjectKey:   evidence.ObjectKey,
		FileName:    evidence.FileName,
		ContentType: evidence.ContentType,
		Size:        evidence.Size,
		CreatedAt:   evidence.CreatedAt,
	}
}
