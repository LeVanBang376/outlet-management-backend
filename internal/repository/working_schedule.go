package repository

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
	"time"
)

type WorkingScheduleRepository struct {
	db *sql.DB
}

func NewWorkingScheduleRepository(db *sql.DB) *WorkingScheduleRepository {
	return &WorkingScheduleRepository{
		db: db,
	}
}

func (r *WorkingScheduleRepository) Create(ctx context.Context, tx *sql.Tx, schedule *model.WorkingSchedule) (*dto.WorkingScheduleResponse, error) {
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
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	schedule.ScheduleID = uint(id)

	res := dto.ToWorkingScheduleResponse(schedule, []*dto.EvidenceResponse{}, nil)

	return res, nil
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

func (r *WorkingScheduleRepository) FindByOutletID(
	ctx context.Context,
	tx *sql.Tx,
	outletID uint,
) ([]model.WorkingSchedule, error) {
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
		WHERE outlet_id = ?
		ORDER BY schedule_date DESC
	`

	rows, err := tx.QueryContext(ctx, query, outletID)
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *WorkingScheduleRepository) FindAll(ctx context.Context, tx *sql.Tx) ([]*dto.WorkingScheduleResponse, error) {
	query := `
	SELECT
		ws.schedule_id,
		ws.outlet_id,
		o.name AS outlet_name,
		ws.sales_id,
		ws.address,
		ws.schedule_date,
		ws.current_stage,
		ws.expected_stage,
		ws.note,
		ws.sync_status,
		ws.created_at,
		ws.updated_at
	FROM working_schedules ws
	INNER JOIN outlets o
		ON ws.outlet_id = o.outlet_id`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	schedules := make([]*dto.WorkingScheduleResponse, 0)

	for rows.Next() {
		var schedule dto.WorkingScheduleResponse

		err := rows.Scan(
			&schedule.ScheduleID,
			&schedule.OutletID,
			&schedule.OutletName,
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

		schedules = append(schedules, &schedule)
	}

	return schedules, rows.Err()
}

func (r *WorkingScheduleRepository) Update(ctx context.Context, tx *sql.Tx, schedule *model.WorkingSchedule) (*dto.WorkingScheduleResponse, error) {
	query := `
		UPDATE working_schedules
		SET
			sales_id = ?,
			address = ?,
			schedule_date = ?,
			expected_stage = ?,
			sync_status = ?,
			note = ?,
			updated_at = ?
		WHERE schedule_id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		schedule.SalesID,
		schedule.Address,
		schedule.ScheduleDate,
		schedule.ExpectedStage,
		schedule.SyncStatus,
		schedule.Note,
		schedule.UpdatedAt,
		schedule.ScheduleID,
	)

	if err != nil {
		return nil, err
	}

	res := dto.ToWorkingScheduleResponse(schedule, []*dto.EvidenceResponse{}, nil)
	return res, nil
}

func (r *WorkingScheduleRepository) Delete(ctx context.Context, tx *sql.Tx, id uint) error {
	query := `
		DELETE FROM working_schedules
		WHERE schedule_id = ?
	`

	_, err := tx.ExecContext(ctx, query, id)

	return err
}

func (r *WorkingScheduleRepository) FindBySalesOutletAndDate(
	ctx context.Context,
	tx *sql.Tx,
	salesID uint,
	outletID uint,
	scheduleDate time.Time,
) (*model.WorkingSchedule, error) {
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
		WHERE sales_id = ?
			AND outlet_id = ?
			AND schedule_date = ?
		LIMIT 1
	`

	var schedule model.WorkingSchedule

	err := tx.QueryRowContext(
		ctx,
		query,
		salesID,
		outletID,
		scheduleDate,
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

func (r *WorkingScheduleRepository) Upsert(
	ctx context.Context,
	tx *sql.Tx,
	schedule *model.WorkingSchedule,
) (*dto.WorkingScheduleResponse, error) {
	existing, err := r.FindBySalesOutletAndDate(
		ctx,
		tx,
		schedule.SalesID,
		schedule.OutletID,
		schedule.ScheduleDate,
	)

	if err != nil {
		return nil, err
	}

	if existing == nil {
		return r.Create(ctx, tx, schedule)
	}

	// Giữ nguyên ID và CreatedAt của bản ghi cũ
	schedule.ScheduleID = existing.ScheduleID
	schedule.CreatedAt = existing.CreatedAt

	return r.Update(ctx, tx, schedule)
}

func (r *WorkingScheduleRepository) UpdateSyncStatus(
	ctx context.Context,
	id uint,
	status string,
) error {
	query := `
		UPDATE working_schedules
		SET
			sync_status = ?,
			updated_at = ?
		WHERE schedule_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		status,
		time.Now(),
		id,
	)

	return err
}
