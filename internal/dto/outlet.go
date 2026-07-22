package dto

import (
	"magnolia-test-backend/internal/constants"
	customerrors "magnolia-test-backend/internal/custom-errors"
	"magnolia-test-backend/internal/model"
	"time"
)

type CreateOutletRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Channel string `json:"channel"`
	Tier    string `json:"tier"`
	SalesID uint   `json:"sales_id"`
	Stage   string `json:"stage"`
	Note    string `json:"note"`

	HasWorkingSchedule bool   `json:"has_working_schedule"`
	ScheduleDate       string `json:"schedule_date"`
	ExpectedStage      string `json:"expected_stage"`
	ScheduleNote       string `json:"schedule_note"`
}

func (r CreateOutletRequest) Validate() error {
	if r.Name == "" {
		return customerrors.OutletErrNameRequired
	}

	if r.Address == "" {
		return customerrors.OutletErrAddressRequired
	}

	if !constants.ValidChannels[r.Channel] {
		return customerrors.OutletErrInvalidChannel
	}

	if !constants.ValidTiers[r.Tier] {
		return customerrors.OutletErrInvalidTier
	}

	if !constants.ValidStages[r.Stage] {
		return customerrors.OutletErrInvalidStage
	}

	if r.HasWorkingSchedule {
		if _, err := time.Parse("2006-01-02", r.ScheduleDate); err != nil {
			return err
		}
	}

	return nil
}

type UpdateOutletRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Channel string `json:"channel"`
	Tier    string `json:"tier"`
	SalesID uint   `json:"sales_id"`
	Stage   string `json:"stage"`
	Note    string `json:"note"`
}

func (r UpdateOutletRequest) Validate() error {
	if r.Name == "" {
		return customerrors.OutletErrNameRequired
	}

	if r.Address == "" {
		return customerrors.OutletErrAddressRequired
	}

	if !constants.ValidChannels[r.Channel] {
		return customerrors.OutletErrInvalidChannel
	}

	if !constants.ValidTiers[r.Tier] {
		return customerrors.OutletErrInvalidTier
	}

	if !constants.ValidStages[r.Stage] {
		return customerrors.OutletErrInvalidStage
	}

	return nil
}

type OutletResponse struct {
	OutletID         uint                       `json:"outlet_id"`
	Name             string                     `json:"name"`
	Address          string                     `json:"address"`
	Channel          string                     `json:"channel"`
	Tier             string                     `json:"tier"`
	SalesID          uint                       `json:"sales_id"`
	Stage            string                     `json:"stage"`
	Note             string                     `json:"note"`
	WorkingSchedules []*WorkingScheduleResponse `json:"working_schedules"`
}

func ToOutletResponse(
	outlet *model.Outlet,
	schedules []*WorkingScheduleResponse,
) *OutletResponse {
	if outlet == nil {
		return nil
	}

	var note string
	if outlet.Note != nil {
		note = *outlet.Note
	}

	res := &OutletResponse{
		OutletID:         outlet.OutletID,
		Name:             outlet.Name,
		Address:          outlet.Address,
		Channel:          outlet.Channel,
		Tier:             outlet.Tier,
		SalesID:          outlet.SalesID,
		Stage:            outlet.Stage,
		Note:             note,
		WorkingSchedules: make([]*WorkingScheduleResponse, 0, len(schedules)),
	}

	for _, schedule := range schedules {
		res.WorkingSchedules = append(
			res.WorkingSchedules,
			schedule,
		)
	}

	return res
}
