package dto

import (
	"magnolia-test-backend/internal/model"
)

type CreateWorkingScheduleRequest struct {
	OutletID      uint    `json:"outlet_id"`
	SalesID       uint    `json:"sales_id"`
	Address       string  `json:"address"`
	ScheduleDate  string  `json:"schedule_date"`
	CurrentStage  string  `json:"current_stage"`
	ExpectedStage *string `json:"expected_stage"`
	Note          *string `json:"note"`
	SyncStatus    string  `json:"sync_status"`
}

type UpdateWorkingScheduleRequest struct {
	Address       string  `json:"address"`
	ScheduleDate  string  `json:"schedule_date"`
	CurrentStage  string  `json:"current_stage"`
	ExpectedStage *string `json:"expected_stage"`
	Note          *string `json:"note"`
	SyncStatus    string  `json:"sync_status"`
}

type WorkingScheduleResponse struct {
	ScheduleID    uint                `json:"schedule_id"`
	OutletID      uint                `json:"outlet_id"`
	SalesID       uint                `json:"sales_id"`
	Address       string              `json:"address"`
	ScheduleDate  string              `json:"schedule_date"`
	CurrentStage  string              `json:"current_stage"`
	ExpectedStage *string             `json:"expected_stage"`
	Note          *string             `json:"note"`
	SyncStatus    string              `json:"sync_status"`
	Evidences     []*EvidenceResponse `json:"evidences"`
}

func ToWorkingScheduleResponse(
	schedule *model.WorkingSchedule,
	evidences []*EvidenceResponse,
) *WorkingScheduleResponse {
	if schedule == nil {
		return nil
	}

	res := &WorkingScheduleResponse{
		ScheduleID:    schedule.ScheduleID,
		OutletID:      schedule.OutletID,
		SalesID:       schedule.SalesID,
		Address:       schedule.Address,
		ScheduleDate:  schedule.ScheduleDate.Format("2006-01-02"),
		CurrentStage:  schedule.CurrentStage,
		ExpectedStage: schedule.ExpectedStage,
		Note:          schedule.Note,
		SyncStatus:    schedule.SyncStatus,
		Evidences:     make([]*EvidenceResponse, 0, len(evidences)),
	}

	for _, evidence := range evidences {
		res.Evidences = append(
			res.Evidences,
			evidence,
		)
	}

	return res
}
