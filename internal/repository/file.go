package repository

import (
	"context"
	"database/sql"
	"magnolia-test-backend/internal/dto"
	"magnolia-test-backend/internal/model"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(
	ctx context.Context,
	tx *sql.Tx,
	file *model.File,
) (*dto.FileResponse, error) {
	query := `
		INSERT INTO files (
			object_key,
			file_name,
			content_type,
			size,
			created_at
		)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		file.ObjectKey,
		file.FileName,
		file.ContentType,
		file.Size,
		file.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	file.FileID = uint(id)

	return dto.ToFileResponse(file), nil
}

func (r *FileRepository) FindByID(
	ctx context.Context,
	tx *sql.Tx,
	id uint,
) (*model.File, error) {
	query := `
		SELECT
			file_id,
			object_key,
			file_name,
			content_type,
			size,
			created_at
		FROM files
		WHERE file_id = ?
	`

	var file model.File

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&file.FileID,
		&file.ObjectKey,
		&file.FileName,
		&file.ContentType,
		&file.Size,
		&file.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) FindAll(
	ctx context.Context,
	tx *sql.Tx,
) ([]*dto.FileResponse, error) {
	query := `
		SELECT
			file_id,
			object_key,
			file_name,
			content_type,
			size,
			created_at
		FROM files
		ORDER BY file_id DESC
	`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]*dto.FileResponse, 0)

	for rows.Next() {
		var file dto.FileResponse

		err := rows.Scan(
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

		files = append(files, &file)
	}

	return files, rows.Err()
}

func (r *FileRepository) Delete(
	ctx context.Context,
	tx *sql.Tx,
	id uint,
) error {
	query := `
		DELETE FROM files
		WHERE file_id = ?
	`

	_, err := tx.ExecContext(ctx, query, id)

	return err
}
