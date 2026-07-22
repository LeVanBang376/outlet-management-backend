package dto

import (
	"magnolia-test-backend/internal/model"
	"time"
)

type CreateWorkingScheduleRequest struct {
	OutletID      uint    `json:"outlet_id"`
	SalesID       uint    `json:"sales_id"`
	Address       string  `json:"address"`
	ScheduleDate  string  `json:"schedule_date"`
	CurrentStage  string  `json:"current_stage"`
	ExpectedStage *string `json:"expected_stage"`
	Note          *string `json:"note"`
}

type UpdateWorkingScheduleRequest struct {
	SalesID       uint    `json:"sales_id"`
	Address       string  `json:"address"`
	ScheduleDate  string  `json:"schedule_date"`
	ExpectedStage *string `json:"expected_stage"`
	Note          *string `json:"note"`
	FileIDs       []uint  `json:"file_ids"`
}

type WorkingScheduleResponse struct {
	ScheduleID    uint                `json:"schedule_id"`
	OutletID      uint                `json:"outlet_id"`
	OutletName    string              `json:"outlet_name"`
	SalesID       uint                `json:"sales_id"`
	Address       string              `json:"address"`
	ScheduleDate  string              `json:"schedule_date"`
	CurrentStage  string              `json:"current_stage"`
	ExpectedStage *string             `json:"expected_stage"`
	Note          *string             `json:"note"`
	SyncStatus    string              `json:"sync_status"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Outlet        *OutletResponse     `json:"outlet"`
	Evidences     []*EvidenceResponse `json:"evidences"`
}

func ToWorkingScheduleResponse(
	schedule *model.WorkingSchedule,
	evidences []*EvidenceResponse,
	outlet *OutletResponse,
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
		Outlet:        outlet,
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
