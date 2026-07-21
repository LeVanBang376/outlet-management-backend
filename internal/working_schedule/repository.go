package working_schedule

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, schedule *WorkingSchedule) error {
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

	result, err := r.db.ExecContext(
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

func (r *Repository) FindByID(ctx context.Context, id uint) (*WorkingSchedule, error) {
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

	var schedule WorkingSchedule

	err := r.db.QueryRowContext(
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

func (r *Repository) FindAll(ctx context.Context) ([]WorkingSchedule, error) {
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

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	schedules := make([]WorkingSchedule, 0)

	for rows.Next() {
		var schedule WorkingSchedule

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

func (r *Repository) Update(ctx context.Context, schedule *WorkingSchedule) error {
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

	_, err := r.db.ExecContext(
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

func (r *Repository) Delete(ctx context.Context, id uint) error {
	query := `
		DELETE FROM working_schedules
		WHERE schedule_id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
