package repository

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/model"
)

type WorkingScheduleRepository struct {
	db *sql.DB
}

func NewWorkingScheduleRepository(db *sql.DB) *WorkingScheduleRepository {
	return &WorkingScheduleRepository{
		db: db,
	}
}

func (r *WorkingScheduleRepository) Create(ctx context.Context, tx *sql.Tx, schedule *model.WorkingSchedule) error {
	query := `
		INSERT INTO working_schedules (
			outlet_id,
			sales_id,
			address,
			schedule_date,
			current_stage,
			expected_stage,
			note,
			sync_status,
			created_at,
			updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		schedule.OutletID,
		schedule.SalesID,
		schedule.Address,
		schedule.ScheduleDate,
		schedule.CurrentStage,
		schedule.ExpectedStage,
		schedule.Note,
		schedule.SyncStatus,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	schedule.ScheduleID = uint(id)

	return nil
}

func (r *WorkingScheduleRepository) FindByID(ctx context.Context, tx *sql.Tx, id uint) (*model.WorkingSchedule, error) {
	query := `
		SELECT
			schedule_id,
			outlet_id,
			sales_id,
			address,
			schedule_date,
			current_stage,
			expected_stage,
			note,
			sync_status,
			created_at,
			updated_at
		FROM working_schedules
		WHERE schedule_id = ?
	`

	var schedule model.WorkingSchedule

	err := tx.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&schedule.ScheduleID,
		&schedule.OutletID,
		&schedule.SalesID,
		&schedule.Address,
		&schedule.ScheduleDate,
		&schedule.CurrentStage,
		&schedule.ExpectedStage,
		&schedule.Note,
		&schedule.SyncStatus,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &schedule, nil
}

func (r *WorkingScheduleRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]model.WorkingSchedule, error) {
	query := `
		SELECT
			schedule_id,
			outlet_id,
			sales_id,
			address,
			schedule_date,
			current_stage,
			expected_stage,
			note,
			sync_status,
			created_at,
			updated_at
		FROM working_schedules
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	schedules := make([]model.WorkingSchedule, 0)

	for rows.Next() {
		var schedule model.WorkingSchedule

		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.OutletID,
			&schedule.SalesID,
			&schedule.Address,
			&schedule.ScheduleDate,
			&schedule.CurrentStage,
			&schedule.ExpectedStage,
			&schedule.Note,
			&schedule.SyncStatus,
			&schedule.CreatedAt,
			&schedule.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	return schedules, rows.Err()
}

func (r *WorkingScheduleRepository) Update(ctx context.Context, tx *sql.Tx, schedule *model.WorkingSchedule) error {
	query := `
		UPDATE working_schedules
		SET
			address = ?,
			schedule_date = ?,
			current_stage = ?,
			expected_stage = ?,
			note = ?,
			sync_status = ?,
			updated_at = ?
		WHERE schedule_id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		schedule.Address,
		schedule.ScheduleDate,
		schedule.CurrentStage,
		schedule.ExpectedStage,
		schedule.Note,
		schedule.SyncStatus,
		schedule.UpdatedAt,
		schedule.ScheduleID,
	)

	return err
}

func (r *WorkingScheduleRepository) Delete(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `
		DELETE FROM working_schedules
		WHERE schedule_id = ?
	`

	_, err := tx.ExecContext(ctx, query, id)

	return err
}
