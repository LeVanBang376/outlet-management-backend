package repository

import (
	"context"
	"database/sql"

	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
)

type EvidenceRepository struct {
	db *sql.DB
}

func NewEvidenceRepository(db *sql.DB) *EvidenceRepository {
	return &EvidenceRepository{
		db: db,
	}
}

func (r *EvidenceRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	evidence *model.Evidence,
) (*dto.EvidenceResponse, error) {
	query := `
		INSERT INTO evidences (
			schedule_id,
			file_id
		)
		VALUES (?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		evidence.ScheduleID,
		evidence.FileID,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	evidence.EvidenceID = uint(id)

	return dto.ToEvidenceResponse(evidence, nil), nil
}

func (r *EvidenceRepository) FindByID(
	ctx context.Context,
	tx *sql.Tx,
	id uint,
) (*model.Evidence, error) {
	query := `
		SELECT
			evidence_id,
			schedule_id,
			file_id
		FROM evidences
		WHERE evidence_id = ?
	`

	var evidence model.Evidence

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&evidence.EvidenceID,
		&evidence.ScheduleID,
		&evidence.FileID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &evidence, nil
}

func (r *EvidenceRepository) FindByScheduleID(
	ctx context.Context,
	tx *sql.Tx,
	scheduleID uint,
) ([]*dto.EvidenceResponse, error) {
	query := `
		SELECT
			e.evidence_id,
			e.schedule_id,
			e.file_id,

			f.file_id,
			f.object_key,
			f.file_name,
			f.content_type,
			f.size,
			f.created_at
		FROM evidences e
		INNER JOIN files f
			ON e.file_id = f.file_id
		WHERE e.schedule_id = ?
		ORDER BY e.evidence_id ASC
	`

	rows, err := tx.QueryContext(ctx, query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	evidences := make([]*dto.EvidenceResponse, 0)

	for rows.Next() {
		var evidence model.Evidence
		var file model.File

		err := rows.Scan(
			&evidence.EvidenceID,
			&evidence.ScheduleID,
			&evidence.FileID,

			&file.FileID,
			&file.ObjectKey,
			&file.FileName,
			&file.ContentType,
			&file.Size,
			&file.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		evidences = append(
			evidences,
			dto.ToEvidenceResponse(
				&evidence,
				dto.ToFileResponse(&file),
			),
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return evidences, nil
}

func (r *EvidenceRepository) Delete(
	ctx context.Context,
	tx *sql.Tx,
	id uint,
) error {
	query := `
		DELETE FROM evidences
		WHERE evidence_id = ?
	`

	_, err := tx.ExecContext(ctx, query, id)

	return err
}

func (r *EvidenceRepository) DeleteByScheduleID(
	ctx context.Context,
	tx *sql.Tx,
	scheduleID uint,
) error {
	query := `
		DELETE FROM evidences
		WHERE schedule_id = ?
	`

	_, err := tx.ExecContext(
		ctx,
		query,
		scheduleID,
	)

	return err
}
