package dto

import "magnolia-test-backend/internal/model"

type CreateEvidenceRequest struct {
	ScheduleID uint `json:"schedule_id" validate:"required"`
	FileID     uint `json:"file_id" validate:"required"`
}

type EvidenceResponse struct {
	EvidenceID uint          `json:"evidence_id"`
	ScheduleID uint          `json:"schedule_id"`
	FileID     uint          `json:"file_id"`
	File       *FileResponse `json:"file,omitempty"`
}

func ToEvidenceResponse(
	evidence *model.Evidence,
	file *FileResponse,
) *EvidenceResponse {
	if evidence == nil {
		return nil
	}

	return &EvidenceResponse{
		EvidenceID: evidence.EvidenceID,
		ScheduleID: evidence.ScheduleID,
		FileID:     evidence.FileID,
		File:       file,
	}
}
