package evidence

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

func (r *Repository) Create(
	ctx context.Context,
	evidence *Evidence,
) error {

	query := `
		INSERT INTO evidences (
			schedule_id,
			object_key,
			file_name,
			content_type,
			size,
			created_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		evidence.ScheduleID,
		evidence.ObjectKey,
		evidence.FileName,
		evidence.ContentType,
		evidence.Size,
		evidence.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	evidence.EvidenceID = uint(id)

	return nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id uint,
) (*Evidence, error) {

	query := `
		SELECT
			evidence_id,
			schedule_id,
			object_key,
			file_name,
			content_type,
			size,
			created_at
		FROM evidences
		WHERE evidence_id = ?
	`

	var evidence Evidence

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&evidence.EvidenceID,
		&evidence.ScheduleID,
		&evidence.ObjectKey,
		&evidence.FileName,
		&evidence.ContentType,
		&evidence.Size,
		&evidence.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &evidence, nil
}

func (r *Repository) FindByScheduleID(
	ctx context.Context,
	scheduleID uint,
) ([]Evidence, error) {

	query := `
		SELECT
			evidence_id,
			schedule_id,
			object_key,
			file_name,
			content_type,
			size,
			created_at
		FROM evidences
		WHERE schedule_id = ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		scheduleID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	evidences := make([]Evidence, 0)

	for rows.Next() {
		var evidence Evidence

		err := rows.Scan(
			&evidence.EvidenceID,
			&evidence.ScheduleID,
			&evidence.ObjectKey,
			&evidence.FileName,
			&evidence.ContentType,
			&evidence.Size,
			&evidence.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		evidences = append(evidences, evidence)
	}

	return evidences, rows.Err()
}

func (r *Repository) Delete(
	ctx context.Context,
	id uint,
) error {

	query := `
		DELETE FROM evidences
		WHERE evidence_id = ?
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}
