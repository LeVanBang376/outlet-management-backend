package worker

import (
	"context"
	"magnolia-test-backend/internal/constants"
	"magnolia-test-backend/internal/repository"
	"time"
)

type Worker struct {
	WorkingScheduleRepo *repository.WorkingScheduleRepository
}

func NewWorker(workingScheduleRepo *repository.WorkingScheduleRepository) *Worker {
	return &Worker{
		WorkingScheduleRepo: workingScheduleRepo,
	}
}

func (w *Worker) SyncWorkingSchedule(
	ctx context.Context,
	scheduleID uint,
) error {
	time.Sleep(10 * time.Second)

	if err := w.WorkingScheduleRepo.UpdateSyncStatus(
		ctx,
		scheduleID,
		constants.SyncStatusSynced,
	); err != nil {
		return err
	}

	return nil
}
