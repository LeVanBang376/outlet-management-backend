package working_schedule

import "time"

type CreateWorkingScheduleRequest struct {
	OutletID      uint      `json:"outlet_id"`
	SalesID       uint      `json:"sales_id"`
	Address       string    `json:"address"`
	ScheduleDate  time.Time `json:"schedule_date"`
	CurrentStage  string    `json:"current_stage"`
	ExpectedStage *string   `json:"expected_stage"`
	Note          *string   `json:"note"`
	SyncStatus    string    `json:"sync_status"`
}

type UpdateWorkingScheduleRequest struct {
	Address       string    `json:"address"`
	ScheduleDate  time.Time `json:"schedule_date"`
	CurrentStage  string    `json:"current_stage"`
	ExpectedStage *string   `json:"expected_stage"`
	Note          *string   `json:"note"`
	SyncStatus    string    `json:"sync_status"`
}
